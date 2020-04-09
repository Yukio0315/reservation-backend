package entity

import "time"

// User represent user information
type User struct {
	ID           int       `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Name         string    `gorm:"varchar(20);not null" json:"name"`
	Email        string    `gorm:"type:varchar(100);unique_index;not null" json:"email"`
	Password     string    `gorm:"char(60);not null" json:"password"`
	Reservations Reservations
}
