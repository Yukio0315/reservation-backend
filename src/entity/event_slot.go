package entity

import (
	"time"

	"github.com/Yukio0315/reservation-backend/src/util"
)

// EventSlot is google
type EventSlot struct {
	ID                    ID `gorm:"primary_key"`
	EventID               ID
	Start                 time.Time `gorm:"unique_index;not null"`
	ReservationEventSlots ReservationEventSlots
	Reservations          []Reservations `gorm:"many2many:reservation_event_slots"`
}

func (es EventSlot) fullEventSlotID() (id ID) {
	if len(es.ReservationEventSlots) >= util.MAXIMUM {
		return es.EventID
	}
	return 0
}

func (es EventSlot) notFullEventSlotID() (id ID) {
	if len(es.ReservationEventSlots) < util.MAXIMUM {
		return es.EventID
	}
	return 0
}

// EventSlots are slice of EventSlot
type EventSlots []EventSlot

// IsReservable judge whether it is reservable or not
func (ess EventSlots) IsReservable() bool {
	if len(ess) == 0 {
		return false
	}
	for _, es := range ess {
		if es.fullEventSlotID() != 0 {
			return false
		}
	}
	return true
}

func (ess EventSlots) fullEventSlotIDs() (ids []ID) {
	for _, es := range ess {
		if es.fullEventSlotID() != 0 {
			ids = append(ids, es.ID)
		}
	}
	return ids
}

func (ess EventSlots) notFullEventSlotIDs() (ids []ID) {
	for _, es := range ess {
		if es.notFullEventSlotID() != 0 {
			ids = append(ids, es.ID)
		}
	}
	return ids
}

func (ess EventSlots) findEventSlotsByID(ids []ID) (eventSlots EventSlots) {
	for _, id := range ids {
		for _, es := range ess {
			if id == es.ID {
				eventSlots = append(eventSlots, es)
			}
		}
	}
	return eventSlots
}

func (ess EventSlots) removeIDsFromIDs(ids1 []ID, ids2 []ID) (result []ID) {
	duplicateIDs := ess.findDuplicates(ids1, ids2)
	if len(duplicateIDs) == len(ids1) {
		return []ID{}
	}
	if len(duplicateIDs) == 0 {
		return ids1
	}
	return ess.filterNotDuplicate(ids1, duplicateIDs)
}

func (ess EventSlots) findDuplicates(ids1 []ID, ids2 []ID) (duplicateIDs []ID) {
	for _, id1 := range ids1 {
		for _, id2 := range ids2 {
			if id1 == id2 {
				duplicateIDs = append(duplicateIDs, id1)
				continue
			}
		}
	}
	return duplicateIDs
}

func (ess EventSlots) filterNotDuplicate(ids1 []ID, ids2 []ID) []ID {
	for _, id2 := range ids2 {
		for i, id1 := range ids1 {
			if id1 == id2 {
				ids1 = append(ids1[:i], ids1[i+1:]...)
			}
		}
	}
	return ids1
}

func (ess EventSlots) generateDurationsExceptIDs(ids []ID) (ds Durations) {
	eventSlots := ess.findEventSlotsByID(ess.removeIDsFromIDs(ess.notFullEventSlotIDs(), ids))
	if len(eventSlots) == 0 {
		return nil
	}
	d := Duration{}
	minStart, maxStart, tmp := time.Time{}, time.Time{}, time.Time{}
	for _, es := range eventSlots {
		if minStart.IsZero() {
			minStart, maxStart, tmp = es.Start, es.Start, es.Start
		}
		if es.Start.Equal(tmp) {
			maxStart = es.Start
			tmp = es.Start.Add(time.Hour * util.INTERVAL)
			d = Duration{
				Start: minStart,
				End:   maxStart.Add(time.Hour * util.INTERVAL),
			}
		} else {
			ds = append(ds, d)
			minStart = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			maxStart = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			tmp = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		}
	}
	ds = append(ds, d)
	return ds
}
