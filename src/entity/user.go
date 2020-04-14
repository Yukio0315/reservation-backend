package entity

import "time"

// ID is uint
type ID uint

// HashedPassword is slice byte
type HashedPassword []byte

// PlainPassword is string
type PlainPassword string

// UserName is string
type UserName string

// Email is string
type Email string

// User represent user information
type User struct {
	ID           ID `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserName     UserName       `gorm:"varchar(20);not null"`
	Email        Email          `gorm:"type:varchar(100);unique_index;not null"`
	Password     HashedPassword `gorm:"not null"`
	Reservations Reservations
}

// UserID represent user id
type UserID struct {
	ID ID `uri:"id" binding:"required"`
}

// UserPlainPassword represent user password
type UserPlainPassword struct {
	PlainPassword PlainPassword `json:"password" binding:"required"`
}

// UserNewOldPasswords represent old password and new password
type UserNewOldPasswords struct {
	OldPassword PlainPassword `json:"oldPassword" binding:"required"`
	NewPassword PlainPassword `json:"newPassword" binding:"required"`
}

// NewUser represent new user
type NewUser struct {
	UserName UserName
	Email    Email
	Password PlainPassword
}

// UserIDAndPassword represent user id and password
type UserIDAndPassword struct {
	ID       ID             `json:"id" binding:"required"`
	Password HashedPassword `json:"password" binding:"required"`
}

// UserIDAndName represent user id and name
type UserIDAndName struct {
	ID       ID       `json:"id" binding:"required"`
	UserName UserName `json:"userName" binding:"required"`
}

// UserIDAndEmail represent user id and email
type UserIDAndEmail struct {
	ID    ID    `json:"id" binding:"required"`
	Email Email `json:"email" binding:"required"`
}

// UserInput represent user input
type UserInput struct {
	UserName UserName      `json:"userName"`
	Email    Email         `json:"email" binding:"required"`
	Password PlainPassword `json:"password" binding:"required"`
}

// UserInputMailPassword represent user input
type UserInputMailPassword struct {
	Email    Email         `json:"email" binding:"required"`
	Password PlainPassword `json:"password" binding:"required"`
}

// UserProfile is represent user profile
type UserProfile struct {
	CreatedAt           time.Time            `json:"createdAt"`
	UserName            UserName             `json:"userName"`
	Email               Email                `json:"email"`
	ReservationProfiles []ReservationProfile `json:"reservationProfile"`
}
