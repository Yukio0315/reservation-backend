package service

import (
	"strings"
	"time"

	"github.com/Yukio0315/reservation-backend/src/api"
	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/util"
	"github.com/jinzhu/gorm"
	"google.golang.org/api/calendar/v3"
)

// EventService represent event service
type EventService struct {
	us  UserService
	rs  ReservationService
	ess EventSlotService
}

// FindAll find all events and reservations in a 1 month from today.
func (es EventService) FindAll() (events entity.Events, err error) {
	db := db.Init()
	if err = db.Order("start asc").
		Where("start > ? AND start < ?", time.Now().AddDate(0, 0, 1).Format("2006-01-02"), time.Now().AddDate(0, 0, 30).Format("2006-01-02")).
		Preload("EventSlots.ReservationEventSlots").
		Preload("EventSlots").
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
	db := db.Init()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	for _, e := range list.Items {
		if strings.ToLower(e.Summary) == "share" {
			start, end := es.generateStartEndTime(e)
			id, err := es.upSertEvent(tx, start, end)
			if err != nil {
				tx.Rollback()
				return err
			}
			for _, a := range e.Attendees {
				if a.Email == util.EMAIL {
					continue
				}
				userID, err := es.us.FindIDByEmailTx(tx, entity.Email(a.Email))
				if err != nil {
					tx.Rollback()
					return err
				}
				reservationID, err := es.rs.upSertReservation(tx, start, end, userID, e.Id)
				if err != nil {
					tx.Rollback()
					return err
				}
				err = es.ess.upSertEventSlotsAndReservationEventSlots(tx, start, end, id, reservationID)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	return tx.Commit().Error
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

func (es EventService) upSertEvent(tx *gorm.DB, start time.Time, end time.Time) (id entity.ID, err error) {
	event := entity.Event{
		Start: start,
		End:   end,
	}
	storedEvent := entity.Event{}
	if err = tx.FirstOrCreate(&storedEvent, event).Error; err != nil {
		return id, err
	}
	return storedEvent.ID, nil
}
