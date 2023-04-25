package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/AlessioPani/go-booking/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDbRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDbRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int

	stmt := `INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx,
		stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

// InsertRoomRestriction inserts a room restriction into the databases
func (m *postgresDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomId,
		r.ReservationId,
		r.RestrictionId,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomId returns true if availability exists for a room Id and false if no availability exists
func (m *postgresDbRepo) SearchAvailabilityByDatesByRoomId(start time.Time, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT count(id)
			  FROM room_restrictions
		      WHERE room_id = $1 and $2 < end_date and $3 > start_date;
	`
	var numRows int

	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a list of available rooms for the given start and end date
func (m *postgresDbRepo) SearchAvailabilityForAllRooms(start time.Time, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT r.id, r.room_name
			  FROM rooms r
			  WHERE r.id NOT IN (SELECT rr.room_id 
			                     FROM room_restrictions rr 
			                     WHERE $1 < rr.end_date AND $2 > rr.start_date)
	`

	var rooms []models.Room

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomById gets a room by id
func (m *postgresDbRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, room_name, created_at, updated_at
			  FROM rooms 
			  WHERE id = $1
	`

	var room models.Room
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}

// GetUserById returns a user by id
func (m *postgresDbRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, first_name, last_name, email, password, access_level, created_at, updated_at
			  FROM Users
			  WHERE id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUserById updates an user in the database
func (m *postgresDbRepo) UpdateUserById(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user
func (m *postgresDbRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email=$1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil

}

// AllReservations returns a slice of all reservations
func (m *postgresDbRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id,
	                 r.created_at, r.updated_at, rm.id, rm.room_name
			  FROM reservations r
			  LEFT JOIN rooms rm on (r.room_id = rm.id)
			  ORDER BY r.start_date asc
			`

	row, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer row.Close()

	for row.Next() {
		var item models.Reservation
		err := row.Scan(
			&item.ID,
			&item.FirstName,
			&item.LastName,
			&item.Email,
			&item.Phone,
			&item.StartDate,
			&item.EndDate,
			&item.RoomId,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Room.ID,
			&item.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, item)

		if err := row.Err(); err != nil {
			return reservations, err
		}
	}

	return reservations, nil

}

// AllNewReservations returns a slice of new reservations
func (m *postgresDbRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id,
	                 r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
			  FROM reservations r
			  LEFT JOIN rooms rm on (r.room_id = rm.id)
			  WHERE r.processed = 0
			  ORDER BY r.start_date asc
			`

	row, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer row.Close()

	for row.Next() {
		var item models.Reservation
		err := row.Scan(
			&item.ID,
			&item.FirstName,
			&item.LastName,
			&item.Email,
			&item.Phone,
			&item.StartDate,
			&item.EndDate,
			&item.RoomId,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Processed,
			&item.Room.ID,
			&item.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, item)

		if err := row.Err(); err != nil {
			return reservations, err
		}
	}

	return reservations, nil

}
