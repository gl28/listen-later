package models

type Article struct {
	Title string `json:"title"`
	Author string `json:"author"`
	Body string `json:"article_body"`
	SiteName string `json:"sitename"`
	PubDate string `json:"date"`
	OriginalURL string
	AudioURL string
}