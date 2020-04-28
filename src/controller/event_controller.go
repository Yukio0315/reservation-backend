package controller

import (
	"log"
	"net/http"

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
	ec.ts.CreateEventAndEventSlotAndReservationEventSlot() //TODO: cron job
	userID := entity.UserID{}
	if err := c.ShouldBindUri(&userID); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	events, err := ec.es.FindAll()
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	reservations, err := ec.rs.FindByUserID(userID.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if len(events.FullEventIDs()) == 0 && len(reservations) == 0 {
		c.JSON(http.StatusOK, events.MakeDurations())
	} else {
		reservedEventSlotIDs := reservations.FindEventSlotIDsByUserID(userID.ID)
		c.JSON(http.StatusOK, events.GenerateDurations(reservedEventSlotIDs))
	}
}

// Delete delete event and reservations
func (ec EventController) Delete(c *gin.Context) {
	duration := entity.Duration{}
	if err := c.ShouldBindJSON(&duration); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	googleEventIDs, err := ec.ts.DeleteReservationAndEvent(duration)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err = ec.gc.DeleteEvents(googleEventIDs); err != nil {
		log.Print(err)
	}
	c.Status(http.StatusNoContent)
}
