package controller

import (
	"github.com/Yukio0315/reservation-backend/src/api"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	"github.com/gin-gonic/gin"
)

// EventController represent event
type EventController struct {
	gc api.GoogleCalendar
	es service.EventService
	rs service.ReservationService
}

// Show shows reservable events
func (ec EventController) Show(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(400, err)
	}

	events, err := ec.es.FindAll()
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	ds := events.GenerateDurations()
	c.JSON(200, ds)
}
