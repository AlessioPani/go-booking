package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/AlessioPani/go-booking/internal/models"
)

func (m *testDbRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDbRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomId == 2 {
		return 1, errors.New("Error")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the databases
func (m *testDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomId == 1000 {
		return errors.New("Error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomId returns true if availability exists for a room Id and false if no availability exists
func (m *testDbRepo) SearchAvailabilityByDatesByRoomId(start time.Time, end time.Time, roomId int) (bool, error) {
	// set up a test time
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2060-01-01 as start
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("some error")
	}

	// if the start date is after 2049-12-31, then return false,
	// indicating no availability;
	if start.After(t) {
		return false, nil
	}

	// otherwise, we have availability
	return true, nil
}

// SearchAvailabilityForAllRooms returns a list of available rooms for the given start and end date
func (m *testDbRepo) SearchAvailabilityForAllRooms(start time.Time, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	// if the start date is after 2049-12-31, then return empty slice,
	// indicating no rooms are available;
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return rooms, errors.New("some error")
	}

	if start.After(t) {
		return rooms, nil
	}

	// otherwise, put an entry into the slice, indicating that some room is
	// available for search dates
	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}

// GetRoomById gets a room by id
func (m *testDbRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}

func (m *testDbRepo) GetUserById(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDbRepo) UpdateUserById(u models.User) error {
	return nil
}

func (m *testDbRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

// AllReservations returns a slice of all reservations
func (m *testDbRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// AllNewReservations returns a slice of new reservations
func (m *testDbRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// GetReservationById retrieve from the database a reservation by its id
func (m *testDbRepo) GetReservationById(id int) (models.Reservation, error) {
	var reservations models.Reservation

	return reservations, nil
}

// UpdateReservation updates a reservation
func (m *testDbRepo) UpdateReservation(r models.Reservation) error {
	return nil
}

// DeleteReservation deletes a reservation by ID
func (m *testDbRepo) DeleteReservation(r models.Reservation) error {
	return nil
}

// UpdatedProcessedForReservation updates processed value for a reservation
func (m *testDbRepo) UpdatedProcessedForReservation(id int, processed int) error {
	return nil
}
