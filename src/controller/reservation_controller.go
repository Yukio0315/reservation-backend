package controller

import (
	"errors"

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
		c.AbortWithError(400, err)
	}

	duration := entity.Duration{}
	if err := c.ShouldBindJSON(&duration); err != nil {
		c.AbortWithError(400, err)
	}

	u, err := rc.su.FindByID(id.ID)
	if err != nil {
		c.AbortWithError(400, err)
	}
	if !u.Reservations.IsReservable(duration) {
		c.AbortWithError(400, errors.New("invalid durations. Already reserved"))
	}

	googleEventID, err := rc.gc.AddEvent(u.UserToEmailAndName(), duration)
	if err != nil {
		c.AbortWithError(400, err)
	}

	err = rc.sr.CreateModels(id.ID, duration, googleEventID)
	if err != nil {
		c.AbortWithStatus(404)
	}
	c.Status(200)
}

// Cancel controller cancel reservation
func (rc ReservationController) Cancel(c *gin.Context) {
	id := entity.UserID{}
	if err := c.ShouldBindUri(&id); err != nil {
		c.AbortWithError(400, err)
	}

	reservationID := entity.ReservationID{}
	if err := c.ShouldBindJSON(&reservationID); err != nil {
		c.AbortWithError(400, err)
	}

	if err := rc.sr.DeleteReservation(id.ID, reservationID.ReservationID); err != nil {
		c.AbortWithError(400, err)
	}
	c.Status(200)
}
