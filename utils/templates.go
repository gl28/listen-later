package utils

import (
	"html/template"
	"net/http"
)

var templates *template.Template

type IndexContent struct {
	FeedURL string
}

func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

func RunTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}