package entity

import (
	"time"
)

// Reservation represent user's reservation
type Reservation struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint
	EventID   uint
	Start     time.Time
}

// Reservations represent array of Reservation
type Reservations []Reservation

func (rs Reservations) eventIDs() (i []uint) {
	keys := make(map[uint]bool)
	result := []uint{}
	for _, r := range rs {
		if value := keys[r.EventID]; !value {
			keys[r.EventID] = true
			result = append(result, r.EventID)
		}
	}
	return result
}

// GenerateDurations generate Event durations
func (rs Reservations) GenerateDurations() (ds Durations) {
	eventIDs := rs.eventIDs()
	for _, i := range eventIDs {
		minStart := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		maxStart := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		for _, r := range rs {
			if r.EventID == i {
				if minStart.IsZero() {
					minStart = r.Start
					maxStart = r.Start
				}
				if minStart.Before(r.Start) {
					minStart = r.Start
				}
				if maxStart.After(r.Start) {
					maxStart = r.Start
				}
			}
		}
		d := Duration{
			Start: minStart,
			End:   maxStart.Add(time.Hour),
		}
		ds = append(ds, d)
	}
	return ds
}
