package entity

// Permission represent permission
type Permission struct {
	Permission string `gorm:"varchar(10);not null;unique;"`
	Users      []User
}
