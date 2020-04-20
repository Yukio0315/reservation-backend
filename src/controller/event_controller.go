package controller

import (
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	"github.com/gin-gonic/gin"
)

// EventController represent event
type EventController struct {
	rs  service.ReservationService
	es  service.EventService
	ess service.EventSlotService
}

// Show shows reservable events
func (ec EventController) Show(c *gin.Context) {
	ec.es.CreateModels()
	userID := entity.UserID{}
	if err := c.ShouldBindUri(&userID); err != nil {
		c.JSON(400, err)
	}

	events, err := ec.es.FindAll()
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	reservations, err := ec.rs.FindByUserID(userID.ID)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	if len(events.FullEventIDs()) == 0 && len(reservations) == 0 {
		c.JSON(200, events.MakeDurations())
	} else {
		reservedEventSlotIDs := reservations.FindEventSlotIDsByUserID(userID.ID)
		c.JSON(200, events.GenerateDurations(reservedEventSlotIDs))
	}
}
