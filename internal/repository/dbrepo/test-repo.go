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
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	layout := "2006-01-02"
	test_start, _ := time.Parse(layout, "2021-08-11")
	test_end, _ := time.Parse(layout, "2021-08-12")

	if start.After(end) {
		return rooms, errors.New("some error")
	}

	if start.Equal(test_start) && end.Equal(test_end) {
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
	if id > 10001 {
		return room, errors.New("some error")
	}

	return room, nil
}

// GetUserByID returns a user by id
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

// UpdateUser updates a user in the database
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Authenticate authenticates a user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email != "adm@adm.com" {
		return 0, "", errors.New("some error")
	}

	if testPassword != "book" {
		return 0, "", errors.New("some error")
	}

	return 1, "", nil
}

// AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

// AllNewReservations returns a slice of all new reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

// GetReservationByID returns one reservation by ID
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var res models.Reservation
	return res, nil
}

// UpdateUser updates a reservation in the database
func (m *testDBRepo) UpdateReservation(r models.Reservation) error {
	return nil
}

// DeleteReservation deletes one reservation by id
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

// AllRooms returns all rooms
func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	rooms = append(rooms, models.Room{
		ID:       1,
		RoomName: "General's Quarters",
	})
	return rooms, nil
}

// GetRestrictionsForRoomByDate returns restrictions for a room by date range
func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	layout := "2006-01-02"
	startDate1, _ := time.Parse(layout, "2050-01-02")
	endDate1, _ := time.Parse(layout, "2050-01-03")
	startDate2, _ := time.Parse(layout, "2050-01-04")
	endDate2, _ := time.Parse(layout, "2050-01-05")
	restrictions = append(restrictions, models.RoomRestriction{
		ID:            1,
		StartDate:     startDate1,
		EndDate:       endDate1,
		RoomID:        1,
		ReservationID: 1,
		RestrictionID: 1,
	}, models.RoomRestriction{
		ID:            1,
		StartDate:     startDate2,
		EndDate:       endDate2,
		RoomID:        1,
		ReservationID: 0,
		RestrictionID: 2,
	})
	return restrictions, nil
}

// InsertBlockForRoom inserts a room restirction
func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	return nil
}

// DeleteBlockByID deletes a room restriction
func (m *testDBRepo) DeleteBlockByID(id int) error {
	return nil
}
