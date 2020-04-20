package db

import (
	"os"

	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/jinzhu/gorm"

	// setup mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

// Init setup db
func Init() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		panic("failed to load .env file")
	}

	db, err := gorm.Open("mysql", os.Getenv("DB_INFO")+"?parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.User{}, &entity.Reservation{}, &entity.Event{}, &entity.EventSlot{}, &entity.ReservationEventSlot{})

	return db
}
