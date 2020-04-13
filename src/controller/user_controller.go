package controller

import (
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	"github.com/gin-gonic/gin"
)

// UserController type
type UserController struct {
	s service.UserService
}

// Show controlles user information & reservation
func (uc UserController) Show(c *gin.Context) {
	id := entity.ID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(400, err)
	}

	p, err := uc.s.FindUserProfile(id)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.JSON(200, p)
}

// UpdatePassword controller update password
func (uc UserController) UpdatePassword(c *gin.Context) {
	id := entity.ID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(400, err)
	}

	password := entity.Password{}
	if err := c.ShouldBindJSON(&password); err != nil {
		c.JSON(400, err)
	}

	if err := uc.s.UpdatePassword(id, password); err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.Status(200)
}
