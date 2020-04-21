package service

import (
	"strings"
	"time"

	"github.com/Yukio0315/reservation-backend/src/api"
	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/util"
	"github.com/jinzhu/gorm"
)

// TransactionService represent every start of transaction
type TransactionService struct {
	es  EventService
	ess EventSlotService
	us  UserService
	rs  ReservationService
	res ReservationEventSlotService
}

// CreateReservationAndReservationEventSlot insert reservations and reservation_event_slots table
func (ts TransactionService) CreateReservationAndReservationEventSlot(userID entity.ID, duration entity.Duration, googleEventID string) (err error) {
	tx, err := db.BeginTx()
	if err != nil {
		tx.Rollback()
		return err
	}

	reservationID, err := ts.rs.createReservationTx(tx, userID, duration, googleEventID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = ts.createReservationEventSlot(tx, reservationID, duration); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (ts TransactionService) createReservationEventSlot(tx *gorm.DB, reservationID entity.ID, duration entity.Duration) error {
	tmp := duration.Start
	for tmp.Before(duration.End) {
		eventSlotID, err := ts.ess.findIDByStartTx(tx, tmp)
		if err != nil {
			return err
		}
		tmp = tmp.Add(time.Hour * util.INTERVAL)
		if err = ts.res.createModelTx(tx, reservationID, eventSlotID); err != nil {
			return err
		}
	}
	return nil
}

// DeleteReservationAndEvent delete reservations and events by duration using transaction
func (ts TransactionService) DeleteReservationAndEvent(duration entity.Duration) (googleEventIDs []string, err error) {
	tx, err := db.BeginTx()
	if err != nil {
		tx.Rollback()
		return []string{}, err
	}

	googleEventIDs, err = ts.rs.deleteByDurationTx(tx, duration)
	if err != nil {
		tx.Rollback()
		return []string{}, err
	}
	if err = ts.es.deleteTx(tx, duration); err != nil {
		tx.Rollback()
		return []string{}, err
	}
	return googleEventIDs, tx.Commit().Error
}

// CreateEventAndEventSlotAndReservationEventSlot create event models
func (ts TransactionService) CreateEventAndEventSlotAndReservationEventSlot() error {
	gc := api.GoogleCalendar{}
	list, err := gc.GetEventsList()
	if err != nil {
		return err
	}
	tx, err := db.BeginTx()
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, e := range list.Items {
		if strings.ToLower(e.Summary) == "share" {
			start, end := ts.es.generateStartEndTime(e)
			id, err := ts.es.upsertEventTx(tx, start, end)
			if err != nil {
				tx.Rollback()
				return err
			}
			for _, a := range e.Attendees {
				if a.Email == util.EMAIL {
					continue
				}
				userID, err := ts.us.findIDByEmailTx(tx, entity.Email(a.Email))
				if err != nil {
					tx.Rollback()
					return err
				}
				reservationID, err := ts.rs.upsertReservationTx(tx, start, end, userID, e.Id)
				if err != nil {
					tx.Rollback()
					return err
				}
				err = ts.ess.upsertEventSlotsAndReservationEventSlotsTx(tx, start, end, id, reservationID)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	return tx.Commit().Error
}
