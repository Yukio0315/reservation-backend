package service

import (
	"errors"
	"time"

	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/util"
	"github.com/jinzhu/gorm"
	"google.golang.org/api/calendar/v3"
)

// EventService represent event service
type EventService struct{}

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

func (es EventService) upsertEventTx(tx *gorm.DB, start time.Time, end time.Time) (id entity.ID, err error) {
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

func (es EventService) deleteTx(tx *gorm.DB, duration entity.Duration) error {
	events := entity.Events{}
	if err := tx.Where("start >= ? AND end <= ?", duration.Start, duration.End).Find(&events).Error; err != nil {
		return err
	}
	if len(events) == 0 {
		return errors.New("invalid durations")
	}
	if err := tx.Unscoped().
		Where("start >= ? AND end <= ?", duration.Start, duration.End).
		Delete(entity.Event{}).Error; err != nil {
		return err
	}
	return nil
}
