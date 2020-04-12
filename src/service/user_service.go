package service

import (
	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
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
