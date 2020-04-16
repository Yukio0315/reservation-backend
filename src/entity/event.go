package entity

import (
	"fmt"
	"sort"
	"time"

	"github.com/Yukio0315/reservation-backend/src/util"
)

// Event represent available datetime every one hours
type Event struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Start        time.Time `gorm:"unique_index;not null"`
	Reservations []Reservations
}

// Events are array of slot
type Events []Event

// Duration represent start and end time
type Duration struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Durations are slice of events or reservations
type Durations []Duration

func (es Events) filterEvent() (result Events) {
	for _, e := range es {
		if len(e.Reservations) <= util.MAXIMUM {
			result = append(result, e)
		}
	}
	return result
}

func (es Events) sortByStartAsk() (result Events) {
	sort.SliceStable(es, func(i, j int) bool {
		return es[i].Start.Before(es[j].Start)
	})
	return es
}

// GenerateDurations generate durations from sort and filtered events by start
func (es Events) GenerateDurations() (ds Durations) {
	events := es.filterEvent().sortByStartAsk()
	d := Duration{}
	minStart, maxStart, tmp := time.Time{}, time.Time{}, time.Time{}
	for _, e := range events {
		if minStart.IsZero() {
			minStart, maxStart, tmp = e.Start, e.Start, e.Start
		}
		fmt.Println(e.Start, minStart, maxStart, tmp)
		if e.Start.Equal(tmp) {
			maxStart = e.Start
			tmp = e.Start.Add(time.Hour)
			d = Duration{
				Start: minStart,
				End:   maxStart.Add(time.Hour),
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
