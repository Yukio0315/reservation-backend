package service

import (
	"time"

	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/util"
	"github.com/jinzhu/gorm"
)

// EventSlotService is service for EventSlot
type EventSlotService struct{}

// FindByDuration find event by duration
func (ess EventSlotService) FindByDuration(duration entity.Duration) (eventSlots entity.EventSlots, err error) {
	db := db.Init()
	if err = db.Where("start >= ? AND start < ?", duration.Start, duration.End).
		Preload("ReservationEventSlots").
		Find(&eventSlots).Error; err != nil {
		return entity.EventSlots{}, err
	}
	defer db.Close()

	return eventSlots, nil
}

func (ess EventSlotService) upsertEventSlotsAndReservationEventSlotsTx(tx *gorm.DB, start time.Time, end time.Time, eventID entity.ID, reservationID entity.ID) (err error) {
	tmpStart := start
	for tmpStart.Before(end) {
		eventSlot := entity.EventSlot{
			EventID: eventID,
			Start:   tmpStart,
		}
		tmpStart = tmpStart.Add(time.Hour * util.INTERVAL)
		storedEventSlot := entity.EventSlot{}
		if err = tx.FirstOrCreate(&storedEventSlot, eventSlot).Error; err != nil {
			return err
		}
		reservationEventSlot := entity.ReservationEventSlot{
			ReservationID: reservationID,
			EventSlotID:   storedEventSlot.ID,
		}
		storedReservationEventSlot := entity.ReservationEventSlot{}
		if err = tx.FirstOrCreate(&storedReservationEventSlot, reservationEventSlot).Error; err != nil {
			return err
		}
	}
	return nil
}

func (ess EventSlotService) findIDByStartTx(tx *gorm.DB, start time.Time) (id entity.ID, err error) {
	eventSlot := entity.EventSlot{}
	if err = tx.Where("start = ?", start).First(&eventSlot).Error; err != nil {
		return id, err
	}
	return eventSlot.ID, nil
}
