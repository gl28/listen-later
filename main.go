package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gl28/listen-later/models"
	"github.com/gl28/listen-later/routes"
	"github.com/gl28/listen-later/utils"
)

func main() {
	r := routes.Init()

	db := models.Init()
	defer db.Close()

	utils.LoadTemplates("templates/*.html")

	http.Handle("/", r)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if port == "" {
		port = ":8000"
	}
	fmt.Println("Serving on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}