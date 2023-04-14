package repository

import "github.com/AlessioPani/go-booking/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) error
}