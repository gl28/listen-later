package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/eduncan911/podcast"
	"github.com/gl28/listen-later/models"
)

var client = &http.Client{Timeout: 10 * time.Second}

type headResponse struct {
	ContentLength string `json:"Content-Length"`
}

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

		contentLength := article.ContentLength
		
		if contentLength == 0 {
			fmt.Println("it was zero")
			// do a HEAD request to get the audio file's length in bytes
			req, err := http.NewRequest("HEAD", article.AudioURL, nil)
			if err != nil {
				return nil, err
			}
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode == 200 {
				contentLengthString := resp.Header.Values("Content-Length")[0]
				contentLength, _ = strconv.ParseInt(contentLengthString, 10, 64)
	
				// save content length so we don't have to do request next time
				err = article.UpdateContentLength(contentLength)
				if err != nil {
					return nil, err
				}
			} else {
				// if the HEAD request failed, there is likely a permissions issue
				// so skip this item
				continue
			}

		}

		item.AddEnclosure(article.AudioURL, podcast.MP3, contentLength)
		if _, err := p.AddItem(item); err != nil {
			return nil, err
		}

	}
	return &p, nil
}