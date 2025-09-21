package main

import (
	"booking-app/database"
	"booking-app/helper"
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

// UserData represents a booking user
type UserData struct {
	ID              uint `gorm:"primaryKey;autoIncrement"`
	FirstName       string
	LastName        string
	Email           string `gorm:"uniqueIndex"`
	NumberOfTickets uint
}

// Conference represents the conference details
type Conference struct {
	Name             string
	TotalTickets     uint
	RemainingTickets uint
	Bookings         []UserData
	mu               sync.Mutex
	wg               sync.WaitGroup
	DB               *gorm.DB
}

// Global channel for ticket printing
var ticketChannel = make(chan string, 10)

// greetUsers shows intro
func (c *Conference) greetUsers() {
	fmt.Printf("Welcome to %v booking application\n", c.Name)
	fmt.Printf("We have total of %v tickets and %v are still available.\n", c.TotalTickets, c.RemainingTickets)
	fmt.Println("Get your tickets here to attend")
}

// getFirstNames returns all first names
func (c *Conference) getFirstNames() []string {
	firstNames := []string{}
	for _, booking := range c.Bookings {
		firstNames = append(firstNames, booking.FirstName)
	}
	return firstNames
}

// bookTicket saves booking in DB + memory
func (c *Conference) bookTicket(userTickets uint, firstName, lastName, email string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.RemainingTickets -= userTickets

	userData := UserData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		NumberOfTickets: userTickets,
	}

	// Save to database
	result := c.DB.Create(&userData)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			fmt.Printf("Email %v already has a booking.\n", email)
		} else {
			fmt.Printf("Error saving booking: %v\n", result.Error)
		}
		return
	}

	// Keep a copy in memory
	c.Bookings = append(c.Bookings, userData)

	fmt.Printf("Thank you %v %v for booking %v tickets. Confirmation email: %v\n",
		firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", c.RemainingTickets, c.Name)
}

// sendTicket simulates sending ticket via channel
func (c *Conference) sendTicket(userTickets uint, firstName, lastName, email string) {
	defer c.wg.Done()
	time.Sleep(10 * time.Second) // simulate sending delay
	ticket := fmt.Sprintf("#################\nSending ticket:\n %v tickets for %v %v\nto email address %v\n#################", userTickets, firstName, lastName, email)
	ticketChannel <- ticket
}

// getUserInput collects user input safely
func getUserInput() (string, string, string, uint) {
	reader := bufio.NewReader(os.Stdin)
	var userTickets uint

	fmt.Print("Enter your first name: ")
	firstName, _ := reader.ReadString('\n')
	firstName = strings.TrimSpace(firstName)

	fmt.Print("Enter your last name: ")
	lastName, _ := reader.ReadString('\n')
	lastName = strings.TrimSpace(lastName)

	fmt.Print("Enter your email address: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	for {
		fmt.Print("Enter number of tickets: ")
		ticketsStr, _ := reader.ReadString('\n')
		ticketsStr = strings.TrimSpace(ticketsStr)
		n, err := fmt.Sscanf(ticketsStr, "%d", &userTickets)
		if err != nil || n != 1 || userTickets == 0 {
			fmt.Println("Invalid number of tickets. Try again.")
			continue
		}
		break
	}

	return firstName, lastName, email, userTickets
}

func main() {
	// Connect DB
	database.ConnectDatabase()

	// Auto-migrate model
	database.DB.AutoMigrate(&UserData{})

	conference := Conference{
		Name:             "Go Conference",
		TotalTickets:     50,
		RemainingTickets: 50,
		Bookings:         []UserData{},
		DB:               database.DB,
	}

	// Load existing bookings from DB to maintain remaining tickets
	var previousBookings []UserData
	conference.DB.Find(&previousBookings)
	conference.Bookings = previousBookings
	for _, b := range previousBookings {
		conference.RemainingTickets -= b.NumberOfTickets
	}

	// Start ticket printing goroutine
	go func() {
		for msg := range ticketChannel {
			fmt.Println(msg)
		}
	}()

	conference.greetUsers()

	for {
		firstName, lastName, email, userTickets := getUserInput()
		isValidName, isValidEmail, isValidTicketNumber := helper.ValidateUserInput(firstName, lastName, email, userTickets, conference.RemainingTickets)

		if isValidName && isValidEmail && isValidTicketNumber {
			conference.bookTicket(userTickets, firstName, lastName, email)

			conference.wg.Add(1)
			go conference.sendTicket(userTickets, firstName, lastName, email)

			firstNames := conference.getFirstNames()
			fmt.Printf("The first names of bookings are: %v\n", firstNames)

			if conference.RemainingTickets == 0 {
				fmt.Println("Our conference is fully booked out. Come back next year.")
				break
			}
		} else {
			if !isValidName {
				fmt.Println("First name or last name you entered is too short.")
			}
			if !isValidEmail {
				fmt.Println("Email address you entered is invalid (missing @ or .).")
			}
			if !isValidTicketNumber {
				fmt.Println("Number of tickets you entered is invalid or exceeds remaining tickets.")
			}
		}
	}

	conference.wg.Wait()
	close(ticketChannel) // close channel after all tickets sent

	// Show all bookings
	var allBookings []UserData
	conference.DB.Find(&allBookings)
	fmt.Println("\nBookings in Database:")
	for _, b := range allBookings {
		fmt.Printf("%v %v (%v) - %v tickets\n", b.FirstName, b.LastName, b.Email, b.NumberOfTickets)
	}

	// Close DB connection
	sqlDB, err := conference.DB.DB()
	if err == nil {
		defer sqlDB.Close()
	}
}
