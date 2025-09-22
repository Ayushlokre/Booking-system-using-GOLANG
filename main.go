package main

import (
	"booking-app/database"
	"booking-app/helper"
	"booking-app/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type BookingRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Tickets   uint   `json:"tickets"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func main() {
	database.ConnectDatabase()
	database.DB.AutoMigrate(&models.Conference{}, &models.UserData{})

	conference := models.Conference{
		Name:             "Go Conference",
		TotalTickets:     50,
		RemainingTickets: 50,
	}

	log.Printf("Checking/Creating conference...")
	result := database.DB.FirstOrCreate(&conference, models.Conference{Name: conference.Name})
	if result.Error != nil {
		log.Printf("Error creating/finding conference: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Created new conference: %s (ID: %d)", conference.Name, conference.ID)
	} else {
		log.Printf("Found existing conference: %s (ID: %d)", conference.Name, conference.ID)
	}

	http.HandleFunc("/api/book", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req BookingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSONError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var currentConference models.Conference
		if err := database.DB.First(&currentConference, conference.ID).Error; err != nil {
			sendJSONError(w, "Conference not found", http.StatusInternalServerError)
			return
		}

		isValidName, isValidEmail, isValidTicketNumber := helper.ValidateUserInput(
			req.FirstName, req.LastName, req.Email, req.Tickets, currentConference.RemainingTickets)

		if !isValidName {
			sendJSONError(w, "Invalid name", http.StatusBadRequest)
			return
		}
		if !isValidEmail {
			sendJSONError(w, "Invalid email", http.StatusBadRequest)
			return
		}
		if !isValidTicketNumber {
			sendJSONError(w, "Invalid number of tickets", http.StatusBadRequest)
			return
		}

		success := bookTicket(&currentConference, req.Tickets, req.FirstName, req.LastName, req.Email)
		if !success {
			sendJSONError(w, "Booking failed. Not enough tickets or duplicate email.", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"message": fmt.Sprintf("Thank you %s %s for booking %d tickets!", req.FirstName, req.LastName, req.Tickets),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/api/bookings", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodGet {
			sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var allBookings []models.UserData
		database.DB.Preload("Conference").Find(&allBookings)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allBookings)
	})

	http.HandleFunc("/api/conference", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodGet {
			sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var currentConference models.Conference
		database.DB.First(&currentConference, conference.ID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentConference)
	})

	http.HandleFunc("/api/debug", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		var conf models.Conference
		err := database.DB.First(&conf).Error

		var bookings []models.UserData
		bookingErr := database.DB.Find(&bookings).Error

		debugInfo := map[string]interface{}{
			"conference_error": fmt.Sprintf("%v", err),
			"conference":       conf,
			"bookings_error":   fmt.Sprintf("%v", bookingErr),
			"bookings_count":   len(bookings),
			"bookings":         bookings,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(debugInfo)
	})

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func bookTicket(c *models.Conference, userTickets uint, firstName, lastName, email string) bool {
	log.Printf("Starting bookTicket function")
	log.Printf("Conference ID: %d, Tickets requested: %d", c.ID, userTickets)
	log.Printf("User: %s %s, Email: %s", firstName, lastName, email)

	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Printf("Transaction begin error: %v", tx.Error)
		return false
	}
	log.Printf("Transaction started successfully")

	var currentConference models.Conference
	if err := tx.First(&currentConference, c.ID).Error; err != nil {
		log.Printf("Conference not found (ID: %d): %v", c.ID, err)
		tx.Rollback()
		return false
	}
	log.Printf("Conference found: %s (ID: %d)", currentConference.Name, currentConference.ID)
	log.Printf("Total tickets: %d, Remaining: %d", currentConference.TotalTickets, currentConference.RemainingTickets)

	if userTickets > currentConference.RemainingTickets {
		log.Printf("Not enough tickets. Requested: %d, Available: %d", userTickets, currentConference.RemainingTickets)
		tx.Rollback()
		return false
	}
	log.Printf("Ticket availability check passed")

	userData := models.UserData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		NumberOfTickets: userTickets,
		ConferenceID:    currentConference.ID,
	}

	log.Printf("Attempting to create booking record...")
	if err := tx.Create(&userData).Error; err != nil {
		log.Printf("Failed to create booking: %v", err)
		log.Printf("This might be a duplicate email or database constraint issue")
		tx.Rollback()
		return false
	}
	log.Printf("Booking record created successfully")

	newRemainingTickets := currentConference.RemainingTickets - userTickets
	log.Printf("Updating remaining tickets from %d to %d", currentConference.RemainingTickets, newRemainingTickets)

	if err := tx.Model(&models.Conference{}).Where("id = ?", currentConference.ID).
		Update("remaining_tickets", newRemainingTickets).Error; err != nil {
		log.Printf("Failed to update tickets: %v", err)
		tx.Rollback()
		return false
	}
	log.Printf("Ticket count updated successfully")

	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit error: %v", err)
		return false
	}
	log.Printf("Transaction committed successfully")

	log.Printf("Booking successful for %s %s - %d tickets", firstName, lastName, userTickets)
	return true
}
