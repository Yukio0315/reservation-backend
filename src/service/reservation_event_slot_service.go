package service

import (
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/jinzhu/gorm"
)

// ReservationEventSlotService represent ReservationEventSlot
type ReservationEventSlotService struct{}

// CreateModel create reservationEventSlot model
func (res ReservationEventSlotService) createModelTx(tx *gorm.DB, reservationID entity.ID, eventSlotID entity.ID) error {
	reservationEventSlot := entity.ReservationEventSlot{
		ReservationID: reservationID,
		EventSlotID:   eventSlotID,
	}
	if err := tx.Create(&reservationEventSlot).Error; err != nil {
		return err
	}
	return nil
}
