
---

# Conference Booking System

A **full-stack web application** designed to manage conference ticket bookings efficiently. This project demonstrates the integration of **Next.js (frontend), Go (backend), and PostgreSQL (database)** to handle user interactions, validate input, and maintain ticket availability in real-time.

This system is ideal for small-to-medium conferences and events where organizers want to track ticket sales and ensure proper management of user bookings.

---

## Table of Contents

1. [Features](#features)
2. [Architecture](#architecture)
3. [Technologies Used](#technologies-used)
4. [Project Structure](#project-structure)
5. [Backend Implementation](#backend-implementation)
6. [Frontend Implementation](#frontend-implementation)
7. [Database Design](#database-design)
8. [API Endpoints](#api-endpoints)
9. [Installation & Setup](#installation--setup)
10. [Usage](#usage)
11. [Contributing](#contributing)
12. [License](#license)

---

## Features

* **Real-Time Booking:** Users can book tickets for a conference instantly.
* **Multiple Bookings:** Users can book multiple times with the same email.
* **Input Validation:** Ensures names, emails, and ticket quantities are valid.
* **Ticket Availability Management:** The system updates remaining tickets automatically.
* **Booking Summary:** View all user bookings and conference details.
* **Robust Error Handling:** Handles duplicate emails, overbooking, and invalid requests gracefully.
* **Debugging Endpoint:** Provides a debug view to see all bookings and conference info for developers.

---

## Architecture

The system follows a **client-server architecture**:

1. **Frontend (Next.js):**

   * Handles user interface, input forms, and displays booking/conference data.
   * Communicates with the backend via REST API endpoints.

2. **Backend (Go):**

   * Provides API endpoints for booking tickets, fetching bookings, and conference details.
   * Uses **GORM** ORM to manage database interactions with PostgreSQL.
   * Implements **transactional booking** to ensure ticket count accuracy.

3. **Database (PostgreSQL):**

   * Stores conference details and user bookings.
   * Supports multiple bookings per user and enforces constraints like ticket count validation.

---

## Technologies Used

| Layer      | Technology                  | Description                                                  |
| ---------- | --------------------------- | ------------------------------------------------------------ |
| Frontend   | Next.js, TypeScript, React  | User interface, dynamic forms, API integration               |
| Backend    | Go (Golang), net/http, GORM | Handles requests, manages bookings and database transactions |
| Database   | PostgreSQL                  | Stores conferences and booking data reliably                 |
| ORM        | GORM                        | Simplifies database operations in Go                         |
| Styling    | Tailwind CSS (optional)     | Provides a modern UI look                                    |
| Versioning | Git, GitHub                 | Source control and collaboration                             |

---

## Project Structure

```
booking-app/
│
├─ booking-backend/         # Go backend
│  ├─ main.go               # Main server and API logic
│  ├─ models/               # Database models (Conference, UserData)
│  ├─ database/             # DB connection and initialization
│  └─ helper/               # Validation helper functions
│
├─ booking-frontend/        # Next.js frontend
│  ├─ src/
│  │   ├─ app/             # Pages and components
│  │   └─ components/      # UI elements (Table, Input, Button)
│  ├─ public/              # Static assets
│  └─ tsconfig.json        # TypeScript configuration
│
└─ README.md
```

---

## Backend Implementation

1. **Booking Logic:**

   * Users submit booking requests via `/api/book`.
   * Backend validates names, emails, and requested ticket numbers.
   * A database transaction ensures tickets are deducted atomically.
   * Multiple bookings per user are allowed by removing the `unique` constraint on email.

2. **Concurrency Handling:**

   * The Go backend uses **database transactions** to avoid race conditions on ticket availability.
   * Remaining tickets are always checked before booking to prevent overbooking.

3. **API Routes:**

   * `/api/book` – POST, books tickets.
   * `/api/bookings` – GET, fetch all bookings.
   * `/api/conference` – GET, fetch conference info.
   * `/api/debug` – GET, for development/debugging purposes.

---

## Frontend Implementation

1. **Booking Form:**

   * Users provide first name, last name, email, and number of tickets.
   * Client-side validation ensures the inputs are correctly formatted.

2. **Booking Response:**

   * Displays success message after booking.
   * Fetches and updates booking data in real-time without refreshing the page.

3. **Booking Table:**

   * Shows all bookings dynamically using the `/api/bookings` endpoint.
   * Displays user details and number of tickets booked.

---

## Database Design

### Conference Table

| Column             | Type      | Description              |
| ------------------ | --------- | ------------------------ |
| id                 | bigint    | Primary key              |
| name               | text      | Conference name          |
| total\_tickets     | uint      | Total tickets available  |
| remaining\_tickets | uint      | Tickets left for booking |
| created\_at        | timestamp | Record creation time     |
| updated\_at        | timestamp | Last update time         |

### UserData Table

| Column              | Type      | Description                            |
| ------------------- | --------- | -------------------------------------- |
| id                  | bigint    | Primary key                            |
| first\_name         | text      | User first name                        |
| last\_name          | text      | User last name                         |
| email               | text      | User email (multiple bookings allowed) |
| number\_of\_tickets | uint      | Tickets booked                         |
| conference\_id      | bigint    | Foreign key referencing Conference     |
| created\_at         | timestamp | Record creation time                   |
| updated\_at         | timestamp | Last update time                       |

---

## API Endpoints

### Book Tickets

* **URL:** `/api/book`
* **Method:** POST
* **Request Body:**

```json
{
  "firstName": "John",
  "lastName": "Doe",
  "email": "john@example.com",
  "tickets": 2
}
```

* **Response Body:**

```json
{
  "message": "Thank you John Doe for booking 2 tickets!"
}
```

### Get All Bookings

* **URL:** `/api/bookings`
* **Method:** GET
* **Response Body:**

```json
[
  {
    "ID": 1,
    "FirstName": "John",
    "LastName": "Doe",
    "Email": "john@example.com",
    "NumberOfTickets": 2,
    "Conference": { "Name": "Go Conference" }
  }
]
```

### Get Conference Details

* **URL:** `/api/conference`
* **Method:** GET
* **Response Body:**

```json
{
  "ID": 1,
  "Name": "Go Conference",
  "TotalTickets": 50,
  "RemainingTickets": 48
}
```

### Debug Endpoint

* **URL:** `/api/debug`
* **Method:** GET
* Provides complete debug info including conference data and bookings count.

---

## Installation & Setup

### Prerequisites

* Node.js (v18+)
* Go (v1.20+)
* PostgreSQL

### Backend Setup

```bash
cd booking-app/booking-backend
# Configure DB connection in database package
go run main.go
```

### Frontend Setup

```bash
cd booking-app/booking-frontend
npm install
npm run dev
```

Open in browser: `http://localhost:3000`

---

## Usage

1. Open the frontend URL in a browser.
2. Fill out the booking form with your details.
3. Submit booking – success message will appear.
4. View all bookings on the booking table.
5. Multiple bookings are allowed with the same email.

---

## Contributing

* Fork the repository
* Create a new branch (`git checkout -b feature-name`)
* Make your changes and commit (`git commit -m "Add feature"`)
* Push to the branch (`git push origin feature-name`)
* Open a Pull Request

---



