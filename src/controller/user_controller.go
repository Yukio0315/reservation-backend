package controller

import (
	"github.com/Yukio0315/reservation-backend/src/service"
)

// UserController type
type UserController struct {
	s service.UserService
}

// Index controlles user
// func (uc UserController) Index(c *gin.Context) {
// 	p, err := uc.s.GetAll()

// 	if err != nil {
// 		c.AbortWithStatus(404)
// 	} else {
// 		c.JSON(200, p)
// 	}
// }
