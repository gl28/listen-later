# Listen Later

Allows you to save articles for later by converting them to audio and adding them to your own private podcast feed.

Currently live at https://listen-l8r.herokuapp.com

## How it works

The application is built from:
* the main server, written in Go, which reads/writes to the database and calls the other APIs
* a Postgres database which stores user info and details about saved articles
* two AWS Lambda functions, one for extracting text from articles and one for converting the text to audio
* AWS API Gateway, which allows the application to call the Lambdas
* an AWS S3 bucket for storing audio files

When you submit a URL to the application, it passes that URL to the first Lambda function, which extracts the main content from the article and passes that content back to the application.

The application then sends that text to the second Lambda function, which uses AWS Polly to convert it to audio. That second Lambda function returns the URL of the new audio, which the application saves in the Postgres database.

Each user has a unique RSS feed URL (e.g. `https://listen-l8r.herokuapp.com/rss/<user-key>/`), which can be subscribed to in any podcast app. When the podcast app makes a request to this URL, the server responds with a listing of all the articles saved by the user. The listing includes a title, description, and audio URL for each article.
