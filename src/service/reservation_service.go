package service

import (
	"time"

	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/util"
	"github.com/jinzhu/gorm"
)

// ReservationService represents reservation service
type ReservationService struct {
	es EventSlotService
}

// CreateModels insert reservations and reservation_event_slots table
func (rs ReservationService) CreateModels(userID entity.ID, duration entity.Duration, googleEventID string) (err error) {
	db := db.Init()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		return err
	}

	reservation := entity.Reservation{
		UserID:        userID,
		Start:         duration.Start,
		End:           duration.End,
		GoogleEventID: googleEventID,
	}
	if err = tx.Create(&reservation).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = rs.createReservationEventSlot(tx, reservation.ID, duration); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// upSertReservation upsert  reservations and reservation_event_slots columns using transaction
func (rs ReservationService) upSertReservation(tx *gorm.DB, start time.Time, end time.Time, userID entity.ID, googleEventID string) (id entity.ID, err error) {
	storedReservation := entity.Reservation{}
	reservation := entity.Reservation{
		UserID:        userID,
		Start:         start,
		End:           end,
		GoogleEventID: googleEventID,
	}
	if err = tx.FirstOrCreate(&storedReservation, reservation).Error; err != nil {
		return id, err
	}
	return storedReservation.ID, nil
}

// FindByUserID can find reservations by user ID
func (rs ReservationService) FindByUserID(userID entity.ID) (reservations entity.Reservations, err error) {
	db := db.Init()

	if err = db.Preload("ReservationEventSlots").
		Where("user_id = ? AND start > ? AND start < ?", userID, time.Now().Format("2006-01-02"), time.Now().AddDate(0, 0, 30).Format("2006-01-02")).
		Find(&reservations).
		Error; err != nil {
		return entity.Reservations{}, err
	}
	defer db.Close()

	return reservations, nil
}

func (rs ReservationService) createReservationEventSlot(tx *gorm.DB, id entity.ID, duration entity.Duration) error {
	tmp := duration.Start
	for tmp.Before(duration.End) {
		eventSlotID, err := rs.es.findIDByStartTX(tx, tmp)
		if err != nil {
			return err
		}
		tmp = tmp.Add(time.Hour * util.INTERVAL)
		reservationEventSlot := entity.ReservationEventSlot{
			ReservationID: id,
			EventSlotID:   eventSlotID,
		}
		if err = tx.Create(&reservationEventSlot).Error; err != nil {
			return err
		}
	}
	return nil
}
