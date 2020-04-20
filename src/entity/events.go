package entity

import (
	"github.com/Yukio0315/reservation-backend/src/util"
)

// Events are array of slot
type Events []Event

// MakeDurations method return durations
func (es Events) MakeDurations() (ds Durations) {
	for _, e := range es {
		d := e.makeDuration()
		ds = append(ds, d)
	}
	return ds
}

// FullEventIDs method returns event ids of full participants
func (es Events) FullEventIDs() (result []ID) {
	ids := []uint{}
	for _, e := range es {
		if e.fullEventID() != 0 {
			ids = append(ids, uint(e.ID))
		}
	}
	uintIds := util.UniqueID(ids)
	for _, id := range uintIds {
		result = append(result, ID(id))
	}
	return result
}

// GenerateDurations generate durations which except full IDs and already reserved IDs
func (es Events) GenerateDurations(eventSlotIDs []ID) (result Durations) {
	for _, e := range es {
		durations := e.EventSlots.generateDurationsExceptIDs(eventSlotIDs)
		result = append(result, durations...)
	}
	return result
}
