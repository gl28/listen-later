package models

import "time"

type Article struct {
	Title string `json:"title"`
	Author string `json:"author"`
	Body string `json:"article_body"`
	SiteName string `json:"sitename"`
	PubDate string `json:"date"`
	OriginalURL string
	AudioURL string
	ID int64
	ContentLength int64
	DateAdded time.Time
}

func SaveNewArticle(userId int, article *Article) error {
	stmt, err := db.Prepare("INSERT INTO articles (user_id, title, audio_url, original_url, author, sitename, pub_date) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		userId,
		article.Title,
		article.AudioURL,
		article.OriginalURL,
		article.Author,
		article.SiteName,
		article.PubDate)
	return err
}

func GetArticlesForUser(userId int) ([]*Article, error) {
	query := `
	SELECT title, audio_url, original_url, author, sitename, pub_date, created_at, content_length, id
	FROM articles
	WHERE user_id = $1
	ORDER BY created_at DESC;
	`
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*Article
	for rows.Next() {
		article := &Article{}
		err := rows.Scan(
			&article.Title,
			&article.AudioURL,
			&article.OriginalURL,
			&article.Author,
			&article.SiteName,
			&article.PubDate,
			&article.DateAdded,
			&article.ContentLength,
			&article.ID)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (article *Article) UpdateContentLength(length int64) error {
	stmt, err := db.Prepare("UPDATE articles SET content_length = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(length, article.ID)
	return err
}
