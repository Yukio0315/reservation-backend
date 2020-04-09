package db

import (
	"os"

	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/jinzhu/gorm"

	// setup mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var (
	db  *gorm.DB
	err error
)

// Init setup db
func Init() {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}

	db, err = gorm.Open("mysql", os.Getenv("DB_INFO")+"?parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.User{}, &entity.Reservation{}, &entity.ReservationSlot{}, &entity.Slot{})
}

// GetDB returns db
func GetDB() *gorm.DB {
	return db
}

// Close the database
func Close() {
	defer db.Close()
}
