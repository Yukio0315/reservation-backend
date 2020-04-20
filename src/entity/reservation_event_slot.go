package entity

// ReservationEventSlot describe reservation start divided into interval time.
type ReservationEventSlot struct {
	ReservationID ID
	EventSlotID   ID
}

// ReservationEventSlots describe slice of ReservationEventSlot
type ReservationEventSlots []ReservationEventSlot
