package models

import "gorm.io/gorm"

// Conference represents the conference
type Conference struct {
	gorm.Model
	Name             string
	TotalTickets     uint
	RemainingTickets uint
	Users            []UserData `gorm:"foreignKey:ConferenceID"`
}

// UserData represents a booking user
type UserData struct {
	gorm.Model
	FirstName       string
	LastName        string
	Email           string `gorm:"uniqueIndex"`
	NumberOfTickets uint
	ConferenceID    uint
}
