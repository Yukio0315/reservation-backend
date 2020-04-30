package service

import (
	"errors"
	"time"

	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/util"
	"github.com/google/uuid"
)

// OneTimeURLService struct
type OneTimeURLService struct{}

// Create create OneTimeURL by userID
func (os OneTimeURLService) Create(userID entity.ID) (o entity.OneTimeURL, err error) {
	db := db.Init()

	o = entity.OneTimeURL{UserID: userID, QueryString: uuid.New().String()}

	if err = db.FirstOrCreate(&o).Error; err != nil {
		return entity.OneTimeURL{}, err
	}
	defer db.Close()

	return o, nil
}

// FindByQueryString create OneTimeURL by userID
func (os OneTimeURLService) FindByQueryString(UUID string) (o entity.OneTimeURL, err error) {
	db := db.Init()

	o = entity.OneTimeURL{}

	if err = db.Where("query_string = ?", UUID).First(&o).Error; err != nil {
		return entity.OneTimeURL{}, err
	}
	defer db.Close()

	return o, nil
}

// FindUserByQueryString create OneTimeURL by userID
func (os OneTimeURLService) FindUserByQueryString(UUID string) (user entity.User, o entity.OneTimeURL, err error) {
	db := db.Init()

	if err = db.Where("query_string = ?", UUID).First(&o).Error; err != nil {
		return entity.User{}, entity.OneTimeURL{}, err
	}
	db.Model(&user).Related(&o)
	defer db.Close()

	return user, o, nil
}

// Delete delete OneTimeURL
func (os OneTimeURLService) Delete(o entity.OneTimeURL) error {
	db := db.Init()

	if err := db.Unscoped().Delete(&o).Error; err != nil {
		return err
	}
	if o.CreatedAt.Add(time.Hour * util.URLLIFETIME).Before(time.Now()) {
		return errors.New("This URL is expired")
	}
	return nil
}
