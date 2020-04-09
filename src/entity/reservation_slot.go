package entity

// ReservationSlot is joining table for reservations and slots
type ReservationSlot struct {
	ReservationID int `gorm:"primary_key;auto_increment:false"`
	SlotID        int `gorm:"primary_key;auto_increment:false"`
}
