package models

import "gorm.io/gorm"

type Conference struct {
	gorm.Model
	Name             string
	TotalTickets     uint
	RemainingTickets uint
	Bookings         []UserData `gorm:"foreignKey:ConferenceID"`
}

type UserData struct {
	gorm.Model
	FirstName       string
	LastName        string
	Email           string // removed `uniqueIndex` to allow multiple bookings per email
	NumberOfTickets uint
	ConferenceID    uint
	Conference      Conference `gorm:"foreignKey:ConferenceID"`
}
