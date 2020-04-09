package entity

import (
	"time"
)

// Reservation represent user's reservation
type Reservation struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    int
	Start     time.Time
	End       time.Time
	Slots     *Slots `gorm:"many2many:reservation_slots;"`
}

// Reservations represent array of Reservation
type Reservations []Reservation
