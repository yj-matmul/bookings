package dbrepo

import (
	"errors"
	"time"

	"github.com/yj-matmul/bookings/internal/models"
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
