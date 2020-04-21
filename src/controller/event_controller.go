package controller

import (
	"log"

	"github.com/Yukio0315/reservation-backend/src/api"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	"github.com/gin-gonic/gin"
)

// EventController represent event
type EventController struct {
	rs service.ReservationService
	es service.EventService
	ts service.TransactionService
	gc api.GoogleCalendar
}

// Show shows reservable events
func (ec EventController) Show(c *gin.Context) {
	ec.ts.CreateEventAndEventSlotAndReservationEventSlot()
	userID := entity.UserID{}
	if err := c.ShouldBindUri(&userID); err != nil {
		c.AbortWithError(400, err)
		return
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

// Delete delete event and reservations
func (ec EventController) Delete(c *gin.Context) {
	duration := entity.Duration{}
	if err := c.ShouldBindJSON(&duration); err != nil {
		c.AbortWithError(400, err)
		return
	}

	googleEventIDs, err := ec.ts.DeleteReservationAndEvent(duration)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err = ec.gc.DeleteEvents(googleEventIDs); err != nil {
		log.Print(400, err)
	}
	c.Status(200)
}
