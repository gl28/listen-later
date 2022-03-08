package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gl28/listen-later/apis"
	"github.com/gl28/listen-later/models"
	"github.com/gl28/listen-later/utils"
	"github.com/gorilla/mux"
)

const rootURL string = "https://listen-l8r.herokuapp.com/"
var ErrUserNotLoggedIn error = errors.New("Could not get user ID from session because user is not logged in.")

func internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	utils.RunTemplate(w, "internal_server_error.html", nil)
	log.Fatal(err)
}

func notFoundError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	utils.RunTemplate(w, "not_found.html", nil)
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromSession(r)
	if err != nil {
		internalServerError(w, err)
		return
	}
	user, err := models.GetUserById(userId)
	if err != nil {
		internalServerError(w, err)
		return
	}
	feedURL := fmt.Sprintf("%srss/%s", rootURL, user.Key)
	articles, err := models.GetArticlesForUser(userId)
	if err !=  nil {
		internalServerError(w, err)
		return
	}
	indexContent := utils.IndexContent{FeedURL: feedURL, Articles: articles}
	utils.RunTemplate(w, "index.html", indexContent)
}

func addArticlePostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.PostForm.Get("articleURL")

	// call first AWS API to extract text content and metadata
	article, err := apis.ExtractContent(url)
	if err != nil {
		internalServerError(w, err)
		return
	}

	// will need user ID to save article after request finishes
	userId, err := getUserIdFromSession(r)
	if err != nil {
		internalServerError(w, err)
		return
	}
	text := article.Body
	request := &apis.AudioConversionRequest{Text: text}

	// send text to AWS API responsible for converting to audio
	audioUrl, err := apis.ConvertToAudio(request)
	if err != nil {
		internalServerError(w, err)
		return
	}
	article.AudioURL = audioUrl

	// if it has no title, set the title to the original URL.
	// some values can be null, but each article must have a
	// title for the RSS feed
	if article.Title == "" {
		article.Title = article.OriginalURL
	}

	err = models.SaveNewArticle(userId, article)

	http.Redirect(w, r, "/?status=success", 303)
}

func rssGetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"]
	user, err := models.GetUserByKey(key)
	if err != nil {
		// most likely url is invalid, so return 404
		notFoundError(w)
		return
	}
	p, err := utils.CreateFeedForUser(user.Id)
	if err != nil {
		internalServerError(w, err)
	}
	w.Header().Set("Content-Type", "application/xml")
	if err := p.Encode(w); err != nil {
		internalServerError(w, err)
		return
	}
}

func Init() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", requireAuthorization(indexGetHandler)).Methods("GET")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	r.HandleFunc("/add", requireAuthorization(addArticlePostHandler)).Methods("POST")
	r.HandleFunc("/rss/{key}", rssGetHandler).Methods("GET")

	return r
}