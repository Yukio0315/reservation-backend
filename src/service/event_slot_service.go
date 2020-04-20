package service

import (
	"time"

	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/util"
	"github.com/jinzhu/gorm"
)

// EventSlotService is service for EventSlot
type EventSlotService struct{}

func (ess EventSlotService) upSertEventSlotsAndReservationEventSlots(tx *gorm.DB, start time.Time, end time.Time, eventID entity.ID, reservationID entity.ID) (err error) {
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

func (ess EventSlotService) findIDByStartTX(tx *gorm.DB, start time.Time) (id entity.ID, err error) {
	eventSlot := entity.EventSlot{}
	if err = tx.Where("start = ?", start).First(&eventSlot).Error; err != nil {
		return id, err
	}
	return eventSlot.ID, nil
}
