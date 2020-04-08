package entity

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type User struct {
	gorm.Model
	Name     string `gorm:"varchar(20);not null" json:"name"`
	Email    string `gorm:"type:varchar(100);unique_index;not null" json:"email"`
	Password string `gorm:"char(60);not null"`
}
