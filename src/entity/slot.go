package entity

import (
	"time"
)

// Slot represent available datetime every one hours
type Slot struct {
	ID           int           `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	Start        time.Time     `gorm:"unique_index;not null"`
	Maximum      int           `gorm:"not null;default 6"`
	CurrentNum   int           `form:"not null;default 0"`
	Reservations *Reservations `gorm:"many2many:reservation_slots;"`
}

// Slots are array of slot
type Slots []Slot
