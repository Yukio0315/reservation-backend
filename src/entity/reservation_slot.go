package entity

// ReservationSlot is joining table for reservations and slots
type ReservationSlot struct {
	ReservationID uint `gorm:"primary_key;auto_increment:false"`
	SlotID        uint `gorm:"primary_key;auto_increment:false"`
}
