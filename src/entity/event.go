package entity

import "time"

// Event represent available datetime every one hours
type Event struct {
	ID         ID `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Start      time.Time `gorm:"unique_index;not null"`
	End        time.Time `gorm:"unique_index;not null"`
	EventSlots EventSlots
}

// fullEventID return id when the event is full
func (e Event) fullEventID() (id ID) {
	ess := e.EventSlots
	if len(ess.fullEventSlotIDs()) != 0 {
		return e.ID
	}
	return id
}

// Duration represent start and end time
type Duration struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Durations are slice of events or reservations
type Durations []Duration

func (e Event) makeDuration() Duration {
	return Duration{
		Start: e.Start,
		End:   e.End,
	}
}
