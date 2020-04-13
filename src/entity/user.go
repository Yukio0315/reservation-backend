package entity

import "time"

// User represent user information
type User struct {
	ID           uint      `gorm:"primary_key" json:"id" binding:"required"`
	CreatedAt    time.Time `json:"createdAt" binding:"required"`
	UpdatedAt    time.Time `json:"updatedAt" binding:"required"`
	UserName     string    `gorm:"varchar(20);not null" json:"username" binding:"required"`
	Email        string    `gorm:"type:varchar(100);unique_index;not null" json:"email" binding:"required"`
	Password     []byte    `gorm:"not null" json:"password" binding:"required"`
	Reservations Reservations
}

// ID represent user id
type ID struct {
	ID uint `uri:"id" binding:"required"`
}

// Password represent user password
type Password struct {
	Password string `json:"password" binding:"required"`
}

// NewUser represent new user
type NewUser struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserAuth represent necessary information for authentication
type UserAuth struct {
	ID       uint   `form:"id" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password []byte `form:"password" json:"password" binding:"required"`
}

// UserInput represent user input
type UserInput struct {
	UserName string `form:"username" json:"username"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// UserProfile is represent user profile
type UserProfile struct {
	CreatedAt           time.Time            `json:"createdAt" binding:"required"`
	UserName            string               `json:"username" binding:"required"`
	Email               string               `json:"email" binding:"required"`
	ReservationProfiles []ReservationProfile `json:"reservationProfile"`
}
