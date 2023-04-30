package models

import (
	"time"
)

// User is the User model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the Room model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction is the Restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is the Reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Processed int
	Room      Room
}

// RoomRestriction is the Room Restriction model
type RoomRestriction struct {
	ID            int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StartDate     time.Time
	EndDate       time.Time
	RoomId        int
	Room          Room
	ReservationId int
	Reservation   Reservation
	RestrictionId int
	Restriction   Restriction
}

// MailData holds an email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
