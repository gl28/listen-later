package main

import (
	"fmt"
	"net/http"

	"github.com/gl28/listen-later/models"
	"github.com/gl28/listen-later/utils"

	"github.com/gorilla/mux"
)

var (
	id int
	email string
	key string
	hash string
)

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.RunTemplate(w, "index.html", nil)
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, "/", 302)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")

	db := models.Init()
	defer db.Close()

	utils.LoadTemplates("templates/*.html")

	http.Handle("/", r)
	fmt.Println("Now serving on localhost:8000...")
	http.ListenAndServe(":8000", nil)
}