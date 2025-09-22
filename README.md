
---

# **Conference Ticket Booking System**

## **Table of Contents**

1. [Overview](#overview)
2. [Tech Stack](#tech-stack)
3. [Database Tables](#database-tables)
4. [Features](#features)
5. [Installation](#installation)

   * [Backend (Go API)](#backend-go-api)
   * [Frontend (Next.js + ShadCN UI)](#frontend-nextjs--shadcn-ui)
6. [How It Works](#how-it-works)
7. [Usage / Booking Flow](#usage--booking-flow)
8. [Notes](#notes)
9. [Project Preview](#project-preview)

---

## **Overview**

This is a **conference ticket booking system** allowing users to book tickets for a conference, view bookings, and manage ticket availability in real time. It combines **Next.js**, **ShadCN UI**, **server actions**, **Go (Golang) API**, and **PostgreSQL** to provide a modern full-stack solution.

---

## **Tech Stack**

| Layer         | Technology           | Purpose                              |
| ------------- | -------------------- | ------------------------------------ |
| Frontend      | Next.js + TypeScript | SSR/CSR and app routing              |
| UI            | ShadCN UI            | Responsive and modern components     |
| Backend       | Go (Golang)          | Business logic and REST API          |
| Database      | PostgreSQL           | Stores conference & booking data     |
| ORM           | GORM                 | Database operations and transactions |
| Communication | RESTful API (JSON)   | Frontend ↔ Backend                   |

---

## **Database Tables**

### **Conference Table**

| Column             | Type      | Description             |
| ------------------ | --------- | ----------------------- |
| id                 | bigint    | Primary Key             |
| name               | text      | Conference Name         |
| total\_tickets     | int       | Total available tickets |
| remaining\_tickets | int       | Remaining tickets       |
| created\_at        | timestamp | Record creation time    |
| updated\_at        | timestamp | Last updated time       |

### **UserData Table**

| Column              | Type      | Description                       |
| ------------------- | --------- | --------------------------------- |
| id                  | bigint    | Primary Key                       |
| first\_name         | text      | User's first name                 |
| last\_name          | text      | User's last name                  |
| email               | text      | User's email (duplicates allowed) |
| number\_of\_tickets | int       | Number of tickets booked          |
| conference\_id      | bigint    | Foreign Key to `Conference`       |
| created\_at         | timestamp | Record creation time              |
| updated\_at         | timestamp | Last updated time                 |

---

## **Features**

* Modern UI with **ShadCN components**
* **Server action integration** for seamless frontend-backend communication
* **Go API** ensures transactional consistency and ticket validation
* **PostgreSQL database** tracks tickets and bookings
* **Multiple bookings allowed per user email**
* Real-time booking updates and ticket availability
* REST API endpoints:

  * `POST /api/book` → Book tickets
  * `GET /api/bookings` → Fetch all bookings
  * `GET /api/conference` → Get conference info
  * `GET /api/debug` → Debug info

---

## **Installation**

### **Backend (Go API)**

1. Install [Go](https://golang.org/dl/)
2. Set up PostgreSQL and create a database, e.g., `booking_db`
3. Update `database.ConnectDatabase()` in `main.go` with your DB credentials
4. Navigate to the backend folder:

```bash
cd booking-app
```

5. Run the backend server:

```bash
go run main.go
```

6. Server will run at `http://localhost:8080`

---

### **Frontend (Next.js + ShadCN UI)**

1. Install [Node.js](https://nodejs.org/)
2. Navigate to the frontend folder:

```bash
cd booking-frontend
```

3. Install dependencies:

```bash
npm install
```

4. Run development server:

```bash
npm run dev
```

5. Open `http://localhost:3000` in your browser

---

## **How It Works**

1. User fills the booking form on the homepage (First Name, Last Name, Email, Tickets).
2. Form triggers a **server action (`bookTickets`)**.
3. Server action sends a POST request to Go API (`/api/book`).
4. Go API validates:

   * Ticket availability
   * User input (name, email, ticket quantity)
   * Creates booking in PostgreSQL
   * Updates remaining tickets
5. Success or error message is returned to the frontend.
6. Users can fetch all bookings or conference info via dedicated endpoints.

---

## **Usage / Booking Flow**

1. Fill in the booking form with your details.
2. Submit → triggers server action → calls Go API → updates DB.
3. Success message appears:

   ```
   Thank you [FirstName] [LastName] for booking [N] tickets!
   ```
4. View all bookings via the bookings table on the frontend.

---

## **Notes**

* The system **supports multiple bookings per email**
* Remaining tickets are **automatically updated** after each booking
* All actions are **transactional** to prevent overbooking
* Frontend uses **ShadCN components** for responsive, modern UI


