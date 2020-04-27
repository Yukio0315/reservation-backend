package entity

import "time"

// OneTimeURL is one time url which is used for password reset
type OneTimeURL struct {
	ID          ID `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      ID
	QueryString string `gorm:"varchar(36);unique_index;not null"`
}

// OneTimeQuery is query string for one time url
type OneTimeQuery struct {
	UUID string `uri:"uuid" binding:"required,uuid4"`
}
