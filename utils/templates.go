package utils

import (
	"html/template"
	"net/http"

	"github.com/gl28/listen-later/models"
)

var templates *template.Template

type IndexContent struct {
	FeedURL string
	Articles []*models.Article
}

func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

func RunTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}