package repository

import (
	"time"

	"github.com/AlessioPani/go-booking/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomId(start time.Time, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start time.Time, end time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
	GetUserById(id int) (models.User, error)
	UpdateUserById(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationById(id int) (models.Reservation, error)
	UpdateReservation(r models.Reservation) error
	DeleteReservation(id int) error
	UpdatedProcessedForReservation(id int, processed int) error
	AllRooms() ([]models.Room, error)
	GetRestrictionsForRoomByDate(roomId int, startDate, endDate time.Time) ([]models.RoomRestriction, error)
	AddBlockForRoom(roomId int, date time.Time) error
	DeleteBlockById(id int) error
}
