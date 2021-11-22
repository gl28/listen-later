package sessions

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

// REMOVE GOTODOTENV FOR PRODUCTION
var err error = godotenv.Load()

var Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
