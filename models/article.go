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

func SaveNewArticle(userId int, article *Article) error {
	stmt, err := db.Prepare("INSERT INTO articles (user_id, title, audio_url, original_url, author, sitename, pub_date) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userId, article.Title, article.AudioURL, article.OriginalURL, article.Author, article.SiteName, article.PubDate)
	return err
}