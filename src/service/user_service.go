package service

import (
	"time"

	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// UserService struct
type UserService struct {
	rs ReservationService
}

// CreateModel create new user
func (s UserService) CreateModel(userName entity.UserName, email entity.Email, password entity.HashedPassword) (u entity.User, err error) {
	db := db.Init()

	u = entity.User{UserName: userName, Email: email, Password: password}

	if err = db.Create(&u).Error; err != nil {
		return entity.User{}, err
	}
	defer db.Close()

	return u, nil
}

// FindByEmail find user by email
func (s UserService) FindByEmail(email entity.Email) (u entity.User, err error) {
	db := db.Init()

	if err = db.Where("email = ?", email).First(&u).Error; err != nil {
		return entity.User{}, err
	}
	defer db.Close()

	return u, nil
}

// FindByID find users by ID
func (s UserService) FindByID(id entity.ID) (u entity.User, err error) {
	db := db.Init()

	if err = db.Preload("Reservations").Where("id = ?", id).First(&u).Error; err != nil {
		return entity.User{}, err
	}
	defer db.Close()
	return u, nil
}

// FindUserProfileByID return user and reservation
func (s UserService) FindUserProfileByID(id entity.ID) (entity.UserProfile, error) {
	db := db.Init()

	u := entity.User{}
	if err := db.Preload("Reservations", "start >= ?", time.Now().Format("2006-01-02"), func(db *gorm.DB) *gorm.DB {
		return db.Order("reservations.start DESC")
	}).Where("id = ?", id).First(&u).Error; err != nil {
		return entity.UserProfile{}, err
	}
	defer db.Close()

	return entity.UserProfile{
		CreatedAt:    u.CreatedAt,
		UserName:     u.UserName,
		Email:        u.Email,
		Permission:   u.Permission,
		Reservations: u.Reservations.FindReservationIDAndDuration(),
	}, nil
}

// UpdatePassword update password from the given id
func (s UserService) UpdatePassword(id entity.ID, plainPassword entity.PlainPassword) error {
	db := db.Init()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	if err != nil {
		return err
	}

	u := entity.User{}
	db.Where("id = ?", id).First(&u)
	if err = db.Model(u).Update("password", hashedPassword).Error; err != nil {
		return err
	}
	defer db.Close()

	return nil
}

// UpdateUserNameByID update userName by ID
func (s UserService) UpdateUserNameByID(id entity.ID, userName entity.UserName) (err error) {
	db := db.Init()

	u := entity.User{}
	db.Where("id = ?", id).First(&u)
	if err = db.Model(u).Update("user_name", userName).Error; err != nil {
		return err
	}
	defer db.Close()

	return nil
}

// UpdateEmailByID update email by ID
func (s UserService) UpdateEmailByID(id entity.ID, email entity.Email) (err error) {
	db := db.Init()

	u := entity.User{}
	db.Where("id = ?", id).First(&u)
	if err = db.Model(u).Update("email", email).Error; err != nil {
		return err
	}
	defer db.Close()

	return nil
}

// DeleteByID delete user by ID
func (s UserService) DeleteByID(id entity.ID) (err error) {
	db := db.Init()

	u := entity.User{}
	if err = db.Where("id = ?", id).Unscoped().Delete(&u).Error; err != nil {
		return err
	}
	defer db.Close()

	return nil
}

func (s UserService) findIDByEmailTx(tx *gorm.DB, email entity.Email) (entity.ID, error) {
	u := entity.User{}
	if err := tx.Where("email = ?", email).First(&u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}
