package service

import (
	"strings"
	"time"

	"github.com/Yukio0315/reservation-backend/src/api"
	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/util"
	"google.golang.org/api/calendar/v3"
)

// EventService represent event service
type EventService struct{}

// FindAll find all events and reservations in a 1 month from today.
func (es EventService) FindAll() (events entity.Events, err error) {
	db := db.Init()
	if err = db.Preload("Reservations").
		Order("start asc").
		Where("start >= ? AND start <= ?", time.Now().Format("2006-01-02"), time.Now().AddDate(0, 1, 0).Format("2006-01-02")).
		Find(&events).Error; err != nil {
		return entity.Events{}, err
	}
	defer db.Close()

	return events, nil
}

// CreateModels create event models
func (es EventService) CreateModels() error {
	gc := api.GoogleCalendar{}
	list, err := gc.GetEventsList()
	if err != nil {
		return err
	}
	for _, e := range list.Items {
		if strings.ToLower(e.Summary) == "share" {
			start, end := es.generateStartEndTime(e)
			daytime := start
			for daytime.Before(end) {
				event := entity.Event{
					Start: daytime,
				}
				daytime = daytime.Add(time.Hour)
				if err = es.updateOrCreateEvent(event); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (es EventService) generateStartEndTime(e *calendar.Event) (start time.Time, end time.Time) {
	if e.Start.DateTime == "" {
		start, _ = time.ParseInLocation("2006-01-02", e.Start.Date, time.Local)
		start = time.Date(start.Year(), start.Month(), start.Day(), util.STARTTIME, 0, 0, 0, time.Local)
		end, _ = time.ParseInLocation("2006-01-02", e.Start.Date, time.Local)
		end = time.Date(end.Year(), end.Month(), end.Day(), util.ENDTIME, 0, 0, 0, time.Local)
	}
	if e.Start.Date == "" {
		start, _ = time.ParseInLocation("2006-01-02T15:04:05Z07:00", e.Start.DateTime, time.Local)
		end, _ = time.ParseInLocation("2006-01-02T15:04:05Z07:00", e.End.DateTime, time.Local)
	}
	return start, end
}

func (es EventService) updateOrCreateEvent(event entity.Event) (err error) {
	db := db.Init()
	eventType := entity.Event{}
	if err = db.FirstOrCreate(&eventType, event).Error; err != nil {
		return err
	}
	return nil
}
