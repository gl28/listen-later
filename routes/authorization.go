package routes

import (
	"net/http"

	"github.com/gl28/listen-later/models"
	"github.com/gl28/listen-later/sessions"
	"github.com/gl28/listen-later/utils"
)

func getUserIdFromSession(r *http.Request) (int, error) {
	session, err := sessions.Store.Get(r, "session")
	if err != nil {
		return 0, err
	}
	userIdUntyped := session.Values["user_id"]

	if userIdUntyped != nil {
		userId := userIdUntyped.(int)
		return userId, nil
	}
	return -1, ErrUserNotLoggedIn
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
		internalServerError(w, err)
		return
	}
	session, err := sessions.Store.Get(r, "session")
	if err != nil {
		internalServerError(w, err)
		return
	}
	session.Values["user_id"] = user.Id
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
		internalServerError(w, err)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func requireAuthorization(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := sessions.Store.Get(r, "session")
		if err != nil {
			internalServerError(w, err)
			return
		}
		_, ok := session.Values["user_id"]

		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
}