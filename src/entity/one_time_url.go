package entity

import "time"

// OneTimeURL is one time url which is used for password reset
type OneTimeURL struct {
	ID         ID `gorm:"primary_key"`
	UserID     ID
	QueryParam string `gorm:"varchar(100);unique_index;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
