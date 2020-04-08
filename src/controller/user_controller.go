package user

import (
	user "github.com/Yukio0315/reservation-backend/src/service"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

func (pc Controller) Index(c *gin.Context) {
	var s user.Service
	p, err := s.GetAll()

	if err != nil {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, p)
	}
}
