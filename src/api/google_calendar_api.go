package api

import (
	"log"
	"time"

	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/joho/godotenv"
	"google.golang.org/api/calendar/v3"
)

// GoogleCalendar represent google calendar
type GoogleCalendar struct{}

func (gc GoogleCalendar) init() (srv *calendar.Service, calendarID string) {
	if err := godotenv.Load(); err != nil {
		panic("failed to load .env file")
	}

	srv, err := calendar.New(googleOauth2Client())
	if err != nil {
		log.Fatal("Failed to connect google calendar client")
	}
	calendarID = "primary"
	return srv, calendarID
}

// GetEventsList shows all events from google calendar
func (gc GoogleCalendar) GetEventsList() (event *calendar.Events, err error) {
	srv, calendarID := gc.init()
	return srv.
		Events.
		List(calendarID).
		TimeMin(time.Now().Format("2006-01-02T15:04:05Z07:00")).
		TimeMax(time.Now().AddDate(0, 0, 30).Format("2006-01-02T15:04:05Z07:00")).
		Do()
}

// AddEvent add reservation for the google calendar
func (gc GoogleCalendar) AddEvent(u entity.EmailAndName, d entity.Duration) (string, error) {
	event := &calendar.Event{
		Summary:  "[Share office] Reservation ( " + string(u.UserName) + " )",
		Location: "1-chōme-24-7 Hashiba Taito City, Tōkyō-to 111-002",
		Start: &calendar.EventDateTime{
			DateTime: d.Start.Format("2006-01-02T15:04:05Z07:00"),
			TimeZone: "Asia/Tokyo",
		},
		End: &calendar.EventDateTime{
			DateTime: d.End.Format("2006-01-02T15:04:05Z07:00"),
			TimeZone: "Asia/Tokyo",
		},
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{Email: string(u.Email)},
		},
	}
	srv, calendarID := gc.init()
	event, err := srv.Events.Insert(calendarID, event).Do()
	return event.Id, err
}

// DeleteEvent delete event from calendar
func (gc GoogleCalendar) DeleteEvent(eventID string) error {
	srv, calendarID := gc.init()
	if err := srv.Events.Delete(calendarID, eventID).Do(); err != nil {
		return err
	}
	return nil
}

// DeleteEvents batch delete events
func (gc GoogleCalendar) DeleteEvents(eventIDs []string) error {
	srv, calendarID := gc.init()
	for _, eventID := range eventIDs {
		if err := srv.Events.Delete(calendarID, eventID).Do(); err != nil {
			return err
		}
	}
	return nil
}
