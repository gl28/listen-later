package sessions

import (
	"os"

	"github.com/gorilla/sessions"
)


var Store = sessions.NewCookieStore([]byte(os.Getenv("LISTEN_LATER_SESSION_KEY")))
