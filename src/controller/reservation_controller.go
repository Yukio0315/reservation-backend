package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/Yukio0315/reservation-backend/src/api"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/service"
	"github.com/gin-gonic/gin"
)

// ReservationController is controller for reservations
type ReservationController struct {
	gc  api.GoogleCalendar
	su  service.UserService
	ess service.EventSlotService
	sr  service.ReservationService
	ts  service.TransactionService
}

// Add method add the reservation and add calendar
func (rc ReservationController) Add(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	duration := entity.Duration{}
	if err := c.ShouldBindJSON(&duration); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	u, err := rc.su.FindByID(id.ID)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if !u.Reservations.IsReservable(duration) {
		c.AbortWithError(http.StatusConflict, errors.New("invalid durations. Already reserved"))
		return
	}
	eventSlots, err := rc.ess.FindByDuration(duration)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if !eventSlots.IsReservable() {
		c.AbortWithError(http.StatusNotFound, errors.New("invalid durations. No event exist"))
		return
	}

	googleEventID, err := rc.gc.AddEvent(u.UserToEmailAndName(), duration)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = rc.ts.CreateReservationAndReservationEventSlot(id.ID, duration, googleEventID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
}

// Cancel controller cancel reservation
func (rc ReservationController) Cancel(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	reservationID := entity.ReservationID{}
	if err := c.ShouldBindJSON(&reservationID); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	googleEventID, err := rc.sr.DeleteReservation(id.ID, reservationID.ReservationID)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err := rc.gc.DeleteEvent(googleEventID); err != nil {
		log.Print(err)
	}
	c.Status(http.StatusNoContent)
}
