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

	db.AutoMigrate(&entity.User{}, &entity.Reservation{}, &entity.Event{}, &entity.EventSlot{}, &entity.ReservationEventSlot{}, &entity.Permission{})
	db.Model(&entity.User{}).AddForeignKey("permission", "permissions(permission)", "RESTRICT", "RESTRICT")
	db.Model(&entity.Reservation{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&entity.EventSlot{}).AddForeignKey("event_id", "events(id)", "CASCADE", "CASCADE")
	db.Model(&entity.ReservationEventSlot{}).
		AddForeignKey("reservation_id", "reservations(id)", "CASCADE", "CASCADE").
		AddForeignKey("event_slot_id", "event_slots(id)", "CASCADE", "CASCADE")

	return db
}

func BeginTx() (*gorm.DB, error) {
	db := Init()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	return tx, tx.Error
}
