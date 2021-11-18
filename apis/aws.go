package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gl28/listen-later/models"
)

type APIResponse struct {
	StatusCode string `json:"statusCode"`
	Body string `json:"body"`
}

var client = &http.Client{Timeout: 10 * time.Second}

func ExtractContent(articleUrl string) (*models.Article, error) {
	/* Calls AWS Lambda responsible for extracting main content
	from article. API will return 400 status if article is invalid */

	API_URL := "https://j4x2f3x9fj.execute-api.us-east-1.amazonaws.com/production"
	body := fmt.Sprintf(`{"article_url":"%s"}`, articleUrl)
	byteStr := []byte(body)
	
	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(byteStr))
	if err != nil {
		return &models.Article{}, err
	}
	req.Header.Add("x-api-key", os.Getenv("AWS_API_KEY_1"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return &models.Article{}, err
	}
	defer resp.Body.Close()

	response := &APIResponse{}
	article := &models.Article{}

	json.NewDecoder(resp.Body).Decode(response)
	json.Unmarshal([]byte(response.Body), article)

	return article, nil
}