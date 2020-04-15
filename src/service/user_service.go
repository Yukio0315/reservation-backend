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
func (s UserService) CreateModel(userName entity.UserName, email entity.Email, password entity.HashedPassword) (entity.User, error) {
	db := db.Init()

	u := entity.User{UserName: userName, Email: email, Password: password}

	if err := db.Create(&u).Error; err != nil {
		return entity.User{}, err
	}
	defer db.Close()

	return u, nil
}

func (s UserService) findByEmail(email entity.Email) (u entity.User, err error) {
	db := db.Init()

	if err = db.Where("email = ?", email).First(&u).Error; err != nil {
		return entity.User{}, err
	}
	defer db.Close()

	return u, nil
}

// FindIDByEmail find a user auth information by email
func (s UserService) FindIDByEmail(email entity.Email) (entity.ID, error) {
	u, err := s.findByEmail(email)
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

// FindIDAndPasswordByEmail find a user auth information by email
func (s UserService) FindIDAndPasswordByEmail(email entity.Email) (entity.UserIDAndPassword, error) {
	u, err := s.findByEmail(email)
	if err != nil {
		return entity.UserIDAndPassword{}, err
	}
	return entity.UserIDAndPassword{
		ID:       u.ID,
		Password: u.Password,
	}, nil
}

// FindEmailAndPasswordByID find a user password by id
func (s UserService) FindEmailAndPasswordByID(id entity.ID) (entity.UserMailAndPassword, error) {
	db := db.Init()

	var u entity.User
	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		return entity.UserMailAndPassword{}, err
	}
	defer db.Close()

	return entity.UserMailAndPassword{
		Email:    u.Email,
		Password: u.Password,
	}, nil
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
	defer db.Close()

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
	defer db.Close()

	return nil
}

// UpdateUserNameByID update userName by ID
func (s UserService) UpdateUserNameByID(id entity.ID, userName entity.UserName) (err error) {
	db := db.Init()

	var u entity.User
	db.Where("id = ?", id).First(&u)
	if err := db.Model(u).Update("user_name", userName).Error; err != nil {
		return err
	}
	defer db.Close()

	return nil
}

// UpdateEmailByID update email by ID
func (s UserService) UpdateEmailByID(id entity.ID, email entity.Email) (err error) {
	db := db.Init()

	var u entity.User
	db.Where("id = ?", id).First(&u)
	if err := db.Model(u).Update("email", email).Error; err != nil {
		return err
	}
	defer db.Close()

	return nil
}

// DeleteByID update email by ID
func (s UserService) DeleteByID(id entity.ID) (err error) {
	db := db.Init()

	var u entity.User
	if err := db.Where("id = ?", id).Unscoped().Delete(&u).Error; err != nil {
		return err
	}
	defer db.Close()

	return nil
}
