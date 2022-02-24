package utils

import (
	"fmt"
	"log"
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

	// get all of user's saved articles from db to add to feed
	articles, err := models.GetArticlesForUser(userId)
	if err != nil {
		return nil, err
	}

	// need user struct to add user's email to feed
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
		
		/* 
			The podcast feed format requires us to include the audio file's length
			in bytes. But the files can take a while to create, so we don't always
			know the length when their info is first added to the database.

			When a user requests their RSS feed we check each audio file's
			contentLength. If contentLength == 0, then we know it hasn't been updated
			yet.

			So we do a HEAD request and save that info in the database for next
			time.
		*/
		if contentLength == 0 {
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
				/*
					If the HEAD request failed, there is likely an S3 permissions issue.
					So skip this item, but log URL for which request failed.
					Continue with other items, because they might not have permissions issues.
				*/
				log.Printf("HEAD request failed when trying to determine content length for: %s", article.AudioURL)
			}

		}

		item.AddEnclosure(article.AudioURL, podcast.MP3, contentLength)
		if _, err := p.AddItem(item); err != nil {
			return nil, err
		}

	}
	return &p, nil
}