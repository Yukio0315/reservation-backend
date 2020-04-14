package service

import (
	"time"

	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// UserService struct
type UserService struct{}

// CreateModel create new user
func (s UserService) CreateModel(username entity.UserName, email entity.Email, password entity.HashedPassword) (entity.User, error) {
	db := db.Init()

	u := entity.User{UserName: username, Email: email, Password: password}

	if err := db.Create(&u).Error; err != nil {
		return entity.User{}, err
	}
	defer db.Close()

	return u, nil
}

// FindByEmail find a user auth information by email
func (s UserService) FindByEmail(email entity.Email) (entity.UserIDAndPassword, error) {
	db := db.Init()

	var u entity.User
	if err := db.Where("email = ?", email).First(&u).Error; err != nil {
		return entity.UserIDAndPassword{}, err
	}
	defer db.Close()

	return entity.UserIDAndPassword{
		ID:       u.ID,
		Password: u.Password,
	}, nil
}

// FindPasswordByID find a user password by id
func (s UserService) FindPasswordByID(id entity.ID) (entity.HashedPassword, error) {
	db := db.Init()

	var u entity.User
	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		return entity.HashedPassword{}, err
	}
	defer db.Close()

	return u.Password, nil
}

// FindUserProfileByID return user and reservation
func (s UserService) FindUserProfileByID(id entity.ID) (entity.UserProfile, error) {
	db := db.Init()

	var u entity.User
	if err := db.Preload("Reservations", "end > ?", time.Now(), func(db *gorm.DB) *gorm.DB {
		return db.Order("reservations.start DESC")
	}).Where("id = ?", id).First(&u).Error; err != nil {
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

// UpdatePassword update password from the given id
func (s UserService) UpdatePassword(id entity.ID, plainPassword entity.PlainPassword) (err error) {
	db := db.Init()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	if err != nil {
		return err
	}

	var u entity.User
	db.Where("id = ?", id).First(&u)
	if err := db.Model(u).Update("password", hashedPassword).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUserNameByID update username by ID
func (s UserService) UpdateUserNameByID(input entity.UserIDAndName) (err error) {
	db := db.Init()

	var u entity.User
	db.Where("id = ?", input.ID).First(&u)
	if err := db.Model(u).Update("user_name", input.UserName).Error; err != nil {
		return err
	}
	return nil
}
