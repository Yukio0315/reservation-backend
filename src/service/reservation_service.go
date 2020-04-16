package service

// ReservationService represents reservation service
type ReservationService struct{}

// func (rs ReservationService)FindByUserID()  {
// 	db := db.Init()
// 	if err = db.Preload("Reservations").
// 		Preload("ReservationSlots").
// 		Where("start >= ?", now).
// 		Order("start asc").
// 		Find(&slots).Error; err != nil {
// 		return entity.Slots{}, err
// 	}
// 	defer db.Close()

// 	return slots, nil
// }
