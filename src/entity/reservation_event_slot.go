package entity

// ReservationEventSlot describe reservation start divided into interval time.
type ReservationEventSlot struct {
	ReservationID ID `gorm:"primary_key;auto_increment:false"`
	EventSlotID   ID `gorm:"primary_key;auto_increment:false"`
}

// ReservationEventSlots describe slice of ReservationEventSlot
type ReservationEventSlots []ReservationEventSlot
