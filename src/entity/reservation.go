package entity

import (
	"time"
)

// Reservation represent user's reservation
type Reservation struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint
	Start     time.Time
	End       time.Time
	Slots     *Slots `gorm:"many2many:reservation_slots;"`
}

// ReservationProfile are for user profile
type ReservationProfile struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Reservations represent array of Reservation
type Reservations []Reservation
