package main

import (
	"fmt"
	"net/http"

	"github.com/gl28/listen-later/apis"
	"github.com/gl28/listen-later/models"
	"github.com/gl28/listen-later/sessions"
	"github.com/gl28/listen-later/utils"

	"github.com/gorilla/mux"
)

func internalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	utils.RunTemplate(w, "internal_server_error.html", nil)
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.RunTemplate(w, "index.html", nil)
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.RunTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	user, err := models.AuthenticateUser(email, password)
	if err == models.ErrInvalidCredentials {
		utils.RunTemplate(w, "login.html", "Invalid username or password.")
		return
	} else if err != nil {
		internalServerError(w)
		return
	}
	session, err := sessions.Store.Get(r, "session")
	if err != nil {
		fmt.Println(err)
		internalServerError(w)
		return
	}
	session.Values["user_id"] = user.Id
	session.Values["email"] = user.Email
	session.Values["user_key"] = user.Key
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.RunTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	err := models.RegisterNewUser(email, password)
	if err == models.ErrUserAlreadyExists {
		utils.RunTemplate(w, "register.html", "A user with that email already exists.")
	} else if err != nil {
		internalServerError(w)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func addArticleGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.RunTemplate(w, "add.html", nil)
}

func addArticlePostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.PostForm.Get("url")

	// call first AWS API to extract text content and metadata
	article, err := apis.ExtractContent(url)
	if err != nil {
		internalServerError(w)
		return
	}

	// send text to AWS API responsible for converting to audio
	
	w.Write([]byte("Success"))
	return
}

func requireAuthorization(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := sessions.Store.Get(r, "session")
		if err != nil {
			internalServerError(w)
			return
		}
		_, ok := session.Values["user_id"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
		}
		handler.ServeHTTP(w, r)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", requireAuthorization(indexGetHandler)).Methods("GET")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	r.HandleFunc("/add", requireAuthorization(addArticleGetHandler)).Methods("GET")
	r.HandleFunc("/add", requireAuthorization(addArticlePostHandler)).Methods("POST")

	db := models.Init()
	defer db.Close()

	utils.LoadTemplates("templates/*.html")

	http.Handle("/", r)
	fmt.Println("Now serving on localhost:8000...")
	http.ListenAndServe(":8000", nil)
}