package api

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/api/calendar/v3"
)

// GoogleCalendar represent google calendar
type GoogleCalendar struct{}

// GetEventsList shows all events from google calendar
func (gc GoogleCalendar) GetEventsList() (*calendar.Events, error) {
	if err := godotenv.Load(); err != nil {
		panic("failed to load .env file")
	}

	srv, err := calendar.New(googleOauth2Client())
	if err != nil {
		log.Fatal("Failed to connect google calendar client")
	}
	address := os.Getenv("GMAIL_ADDRESS")
	return srv.
		Events.
		List(address).
		TimeMin(time.Now().Format("2006-01-02T15:04:05Z07:00")).
		TimeMax(time.Now().AddDate(0, 1, 0).Format("2006-01-02T15:04:05Z07:00")).
		Do()
}
