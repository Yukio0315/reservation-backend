package controller

import (
	"github.com/Yukio0315/reservation-backend/src/api"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	"github.com/gin-gonic/gin"
)

// ReservationController is controller for reservations
type ReservationController struct {
	gc api.GoogleCalendar
	su service.UserService
	sr service.ReservationService
}

// Add method add the reservation and add calendar
func (rc ReservationController) Add(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(400, err)
	}

	duration := entity.Duration{}
	if err := c.ShouldBindJSON(&duration); err != nil {
		c.JSON(400, err)
	}

	u, err := rc.su.FindByID(id.ID)
	if err != nil {
		c.JSON(400, err)
	}

	googleEventID, err := rc.gc.AddEvent(u.UserToEmailAndName(), duration)
	if err != nil {
		c.JSON(400, err)
	}

	err = rc.sr.CreateModels(id.ID, duration, googleEventID)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.Status(200)
}
