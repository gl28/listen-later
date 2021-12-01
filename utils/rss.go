package utils

import (
	"fmt"

	"github.com/eduncan911/podcast"
	"github.com/gl28/listen-later/models"
)

func CreateFeedForUser(userId int) (*podcast.Podcast, error) {
	articles, err := models.GetArticlesForUser(userId)
	if err != nil {
		return nil, err
	}
	user, err := models.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	/* using the date the user's first article was added
		as the feed's publication date, and the date the most recent
		article was added as the feed's last updated date */
	feedUpdatedDate := articles[0].DateAdded
	feedPublishedDate := articles[len(articles)-1].DateAdded
	url := fmt.Sprintf("https://listen-later.herokuapp.com/rss/%s", user.Key)

	p := podcast.New(
		"Listen Later Podcast Feed",
		url,
		"All of the articles you've saved",
		&feedPublishedDate,
		&feedUpdatedDate,
	)
	p.AddAuthor(user.Email, user.Email)

	for _, article := range articles {

		description := fmt.Sprintf(`Site: %s
		Author: %s
		URL: %s
		Publication date: %s`, article.SiteName,
		article.Author, article.OriginalURL, article.PubDate)

		item := podcast.Item{
			Title: article.Title,
			Link: article.AudioURL,
			Description: description,
			PubDate: &article.DateAdded,
		}
		// add content length to enclosure
		item.AddEnclosure(article.AudioURL, podcast.MP3, 1)
		if _, err := p.AddItem(item); err != nil {
			return nil, err
		}

	}
	return &p, nil
}