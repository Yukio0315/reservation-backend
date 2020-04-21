package service

import (
	"time"

	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/jinzhu/gorm"
)

// ReservationService represents reservation service
type ReservationService struct{}

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

// DeleteReservation delete reservations and reservation event slot by userId and reservationId
func (rs ReservationService) DeleteReservation(userID entity.ID, reservationID entity.ID) (string, error) {
	db := db.Init()
	reservation := entity.Reservation{}
	if err := db.Unscoped().
		Where("id = ? AND user_id = ?", reservationID, userID).
		First(&reservation).
		Error; err != nil {
		return "", err
	}
	if err := db.Delete(&reservation).Error; err != nil {
		return "", err
	}
	return reservation.GoogleEventID, nil
}

// upSertReservation upsert  reservations and reservation_event_slots columns using transaction
func (rs ReservationService) upsertReservationTx(tx *gorm.DB, start time.Time, end time.Time, userID entity.ID, googleEventID string) (id entity.ID, err error) {
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

func (rs ReservationService) deleteByDurationTx(tx *gorm.DB, duration entity.Duration) ([]string, error) {
	reservations := entity.Reservations{}
	if err := tx.Where("start >= ? AND end <= ?", duration.Start, duration.End).Find(&reservations).Error; err != nil {
		return []string{}, err
	}
	if err := tx.Unscoped().Where("start >= ? AND end <= ?", duration.Start, duration.End).Delete(entity.Reservation{}).Error; err != nil {
		return []string{}, err
	}
	return reservations.GoogleEventIDs(), nil
}

func (rs ReservationService) createReservationTx(tx *gorm.DB, userID entity.ID, duration entity.Duration, googleEventID string) (id entity.ID, err error) {
	reservation := entity.Reservation{
		UserID:        userID,
		Start:         duration.Start,
		End:           duration.End,
		GoogleEventID: googleEventID,
	}
	if err = tx.Create(&reservation).Error; err != nil {
		return id, err
	}
	return reservation.ID, nil
}
