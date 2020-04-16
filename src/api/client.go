package api

import (
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func googleOauth2Client() *http.Client {
	if err := godotenv.Load(); err != nil {
		panic("failed to load .env file")
	}

	config := oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{"https://mail.google.com/", "https://www.googleapis.com/auth/calendar"},
	}

	expiry, _ := time.Parse("2006-01-02", "2020-04-14")
	token := oauth2.Token{
		AccessToken:  os.Getenv("GOOGLE_ACCESS_TOKEN"),
		TokenType:    "Bearer",
		RefreshToken: os.Getenv("GOOGLE_REFRESH_TOKEN"),
		Expiry:       expiry,
	}

	return config.Client(oauth2.NoContext, &token)

}
