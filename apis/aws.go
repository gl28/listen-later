package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gl28/listen-later/models"
	"github.com/joho/godotenv"
)

type APIResponse struct {
	StatusCode string `json:"statusCode"`
	Body string `json:"body"`
}

type AudioConversionRequest struct {
	Text string `json:"text"`
}


// REMOVE GOTODOTENV FOR PRODUCTION
var err error = godotenv.Load()

var client = &http.Client{Timeout: 10 * time.Second}
var ErrAudioConversionFailed error = errors.New("Server responded, but audio conversion failed.")

func ExtractContent(articleUrl string) (*models.Article, error) {
	/* Calls AWS Lambda responsible for extracting main content and
	metadata from article. API will return 400 status if article is invalid */

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

	article.OriginalURL = articleUrl

	return article, nil
}

func ConvertToAudio(request *AudioConversionRequest) (string, error) {
	/* Calls AWS Lambda which:
		1. Calls Amazon Polly to convert text to audio
		2. Stores resulting audio in S3 bucket
	*/

	API_URL := "https://k022oyyeg5.execute-api.us-east-1.amazonaws.com/production"

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Add("x-api-key", os.Getenv("AWS_API_KEY_2"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	response := &APIResponse{}
	json.NewDecoder(resp.Body).Decode(response)
	url := response.Body

	if response.StatusCode != "200" {
		return "", ErrAudioConversionFailed
	}
	
	return url, nil
}