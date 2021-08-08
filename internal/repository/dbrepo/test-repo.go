package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/yj-matmul/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 10000, then fail; otherwise pass
	if res.RoomID == 10000 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestirction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	// if the room id is 10001, then fail; otherwise pass
	if r.RoomID == 10001 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID == 10000 {
		return false, errors.New("some error")
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	if start.After(end) {
		return rooms, errors.New("some error")
	}

	if start.Equal(end) {
		return rooms, nil
	}

	rooms = append(rooms, models.Room{
		ID:       1,
		RoomName: "General's Quarters",
	})

	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}

// GetUserByID returns a user by id
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
		from users
		where id = $1`

	var u models.User
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates a user in the database
func (m *testDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email= $3, access_level = $4, updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now())

	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(
		&id,
		&hashedPassword,
	)

	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return id, "", errors.New("incorrect password")
	} else if err != nil {
		return id, "", err
	}

	return id, hashedPassword, nil
}
