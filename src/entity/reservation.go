package entity

import (
	"time"

	"github.com/Yukio0315/reservation-backend/src/util"
)

// Reservation represent user's reservation
type Reservation struct {
	ID                    ID `gorm:"primary_key"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	UserID                ID
	GoogleEventID         string `gorm:"not null"`
	Start                 time.Time
	End                   time.Time
	ReservationEventSlots ReservationEventSlots
	EventSlots            []EventSlot `gorm:"many2many:reservation_event_slots;"`
}

func (r Reservation) findEventSlotIDsByReservation() (eventSlotIDs []ID) {
	eventSlotIDsUint := []uint{}
	for _, re := range r.ReservationEventSlots {
		if re.ReservationID == r.ID {
			eventSlotIDsUint = append(eventSlotIDsUint, uint(re.EventSlotID))
		}
	}
	uniqueIDs := util.UniqueID(eventSlotIDsUint)
	for _, id := range uniqueIDs {
		eventSlotIDs = append(eventSlotIDs, ID(id))
	}
	return eventSlotIDs
}

// Reservations represent array of Reservation
type Reservations []Reservation

// ReservationUserAndTime is for add google calendar
type ReservationUserAndTime struct {
	UserName UserName
	Email    Email
	Start    time.Time
	End      time.Time
}

// ReservationIDAndDuration represent reservation information for users
type ReservationIDAndDuration struct {
	ID    ID
	Start time.Time
	End   time.Time
}

// ReservationID represent reservation id
type ReservationID struct {
	ReservationID ID `json:"id" binding:"required"`
}

// FindReservationIDAndDuration return reservation id and duration for user
func (rs Reservations) FindReservationIDAndDuration() (result []ReservationIDAndDuration) {
	for _, r := range rs {
		reservation := ReservationIDAndDuration{
			ID:    r.ID,
			Start: r.Start,
			End:   r.End,
		}
		result = append(result, reservation)
	}
	return result
}

// MakeDurations return durations from reservations
func (rs Reservations) MakeDurations() (ds Durations) {
	for _, r := range rs {
		duration := Duration{
			Start: r.Start,
			End:   r.End,
		}
		ds = append(ds, duration)
	}
	return ds
}

// findIDsByUserID can find ids by user id
func (rs Reservations) findIDsByUserID(userID ID) (ids []ID) {
	for _, r := range rs {
		if r.UserID == userID {
			ids = append(ids, r.ID)
		}
	}
	return ids
}

// FindEventSlotIDsByUserID find eventSlotIDs by userID
func (rs Reservations) FindEventSlotIDsByUserID(userID ID) (eventSlotIDs []ID) {
	reservationIDs := rs.findIDsByUserID(userID)
	for _, r := range rs {
		for _, rid := range reservationIDs {
			if rid == r.ID {
				eventSlotIDs = append(eventSlotIDs, r.findEventSlotIDsByReservation()...)
			}
		}
	}
	return eventSlotIDs
}
