package service

import (
	"time"

	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/jinzhu/gorm"
)

// UserService struct
type UserService struct{}

// CreateModel create new user
func (s UserService) CreateModel(username string, email string, password []byte) (entity.User, error) {
	db := db.Init()

	u := entity.User{UserName: username, Email: email, Password: password}

	if err := db.Create(&u).Error; err != nil {
		return entity.User{}, err
	}
	defer db.Close()

	return u, nil
}

// FindByEmail find a user by email
func (s UserService) FindByEmail(email string) (entity.UserAuth, error) {
	db := db.Init()

	var u entity.User
	if err := db.Where("email = ?", email).First(&u).Error; err != nil {
		return entity.UserAuth{}, err
	}
	defer db.Close()

	return entity.UserAuth{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

// FindUserProfile return user and reservation
func (s UserService) FindUserProfile(id entity.ID) (entity.UserProfile, error) {
	db := db.Init()

	var u entity.User
	if err := db.Preload("Reservations", "end > ?", time.Now(), func(db *gorm.DB) *gorm.DB {
		return db.Order("reservations.start DESC")
	}).Where("id = ?", id.ID).First(&u).Error; err != nil {
		return entity.UserProfile{}, err
	}

	var reservationProfiles []entity.ReservationProfile
	for _, r := range u.Reservations {
		rp := entity.ReservationProfile{
			Start: r.Start,
			End:   r.End,
		}
		reservationProfiles = append(reservationProfiles, rp)
	}

	return entity.UserProfile{
		CreatedAt:           u.CreatedAt,
		UserName:            u.UserName,
		Email:               u.Email,
		ReservationProfiles: reservationProfiles,
	}, nil
}
