package main

import (
	"booking-app/database"
	"booking-app/helper"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Conference struct {
	ID               uint `gorm:"primaryKey;autoIncrement"`
	Name             string
	TotalTickets     uint
	RemainingTickets uint
	Bookings         []UserData
}

type UserData struct {
	ID              uint `gorm:"primaryKey;autoIncrement"`
	FirstName       string
	LastName        string
	Email           string `gorm:"uniqueIndex"`
	NumberOfTickets uint
	ConferenceID    uint
	Conference      Conference `gorm:"foreignKey:ConferenceID"`
}

func (c *Conference) greetUsers() {
	fmt.Printf("Welcome to %v booking application\n", c.Name)
	fmt.Printf("Total tickets: %v, Tickets remaining: %v\n", c.TotalTickets, c.RemainingTickets)
	fmt.Println("Get your tickets here to attend")
	fmt.Println("Type 'exit' as first name to quit anytime.")
}

func (c *Conference) getFirstNames() []string {
	firstNames := []string{}
	for _, booking := range c.Bookings {
		firstNames = append(firstNames, booking.FirstName)
	}
	return firstNames
}

// Book ticket with transaction and rollback
func (c *Conference) bookTicket(userTickets uint, firstName, lastName, email string) bool {
	if userTickets > c.RemainingTickets {
		fmt.Println("Not enough tickets remaining. Try a smaller number.")
		return false
	}

	userData := UserData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		NumberOfTickets: userTickets,
		ConferenceID:    c.ID,
	}

	tx := database.DB.Begin() // start transaction

	// Create user booking
	if err := tx.Create(&userData).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "duplicate key") {
			fmt.Printf("Email %v already has a booking.\n", email)
		} else {
			fmt.Printf("Booking failed: %v\n", err)
		}
		return false
	}

	// Update conference remaining tickets (bulk updatt)
	if err := tx.Model(&Conference{}).Where("id = ?", c.ID).
		Update("remaining_tickets", gorm.Expr("remaining_tickets - ?", userTickets)).Error; err != nil {
		tx.Rollback()
		fmt.Println("Failed to update conference tickets:", err)
		return false
	}

	tx.Commit() // commit transaction

	c.RemainingTickets -= userTickets
	c.Bookings = append(c.Bookings, userData)

	fmt.Printf("Thank you %v %v for booking %v tickets. Confirmation email: %v\n",
		firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", c.RemainingTickets, c.Name)

	return true
}

func sendTicket(userTickets uint, firstName, lastName, email string, delay time.Duration) {
	time.Sleep(delay)
	fmt.Println("#################")
	fmt.Printf("Sending ticket:\n %v tickets for %v %v \nto email: %v\n", userTickets, firstName, lastName, email)
	fmt.Println("#################")
}

func getUserInput() (string, string, string, uint, bool) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your first name: ")
	firstName, _ := reader.ReadString('\n')
	firstName = strings.TrimSpace(firstName)
	if strings.ToLower(firstName) == "exit" {
		return "", "", "", 0, true
	}

	fmt.Print("Enter your last name: ")
	lastName, _ := reader.ReadString('\n')
	lastName = strings.TrimSpace(lastName)

	fmt.Print("Enter your email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter number of tickets: ")
	ticketsStr, _ := reader.ReadString('\n')
	ticketsStr = strings.TrimSpace(ticketsStr)
	userTickets, err := strconv.Atoi(ticketsStr)
	if err != nil || userTickets <= 0 {
		fmt.Println("Invalid number of tickets. Defaulting to 1 ticket.")
		userTickets = 1
	}

	return firstName, lastName, email, uint(userTickets), false
}

func main() {
	// Connect DB and migrate
	database.ConnectDatabase()
	database.DB.AutoMigrate(&Conference{}, &UserData{})

	// Create conference
	conference := Conference{
		Name:             "Go Conference",
		TotalTickets:     50,
		RemainingTickets: 50,
	}
	database.DB.FirstOrCreate(&conference, Conference{Name: conference.Name})

	conference.greetUsers()

	for {
		firstName, lastName, email, userTickets, exit := getUserInput()
		if exit {
			fmt.Println("Exiting the booking application...")
			break
		}

		isValidName, isValidEmail, isValidTicketNumber := helper.ValidateUserInput(firstName, lastName, email, userTickets, conference.RemainingTickets)
		if !isValidName {
			fmt.Println("First name or last name is too short.")
			continue
		}
		if !isValidEmail {
			fmt.Println("Email address is invalid.")
			continue
		}
		if !isValidTicketNumber {
			fmt.Println("Invalid number of tickets.")
			continue
		}

		if conference.bookTicket(userTickets, firstName, lastName, email) {
			sendTicket(userTickets, firstName, lastName, email, 5*time.Second)
			fmt.Printf("The first names of bookings are: %v\n", conference.getFirstNames())

			if conference.RemainingTickets == 0 {
				fmt.Println("All tickets are sold out. Thank you!")
				break
			}
		}
	}

	// Display all bookings using join query
	var allBookings []UserData
	database.DB.Preload("Conference").Find(&allBookings)
	fmt.Println("\nBookings in Database (with Conference):")
	for _, b := range allBookings {
		fmt.Printf("%v %v (%v) - %v tickets, Conference: %v\n",
			b.FirstName, b.LastName, b.Email, b.NumberOfTickets, b.Conference.Name)
	}
}
