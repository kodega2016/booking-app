package dbrepo

import (
	"errors"
	"time"

	"booking-app/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("failed to insert reservation")
	}
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) (int, error) {
	if res.RoomID == 1000 {
		return 0, errors.New("failed to insert room restriction")
	}
	return 1, nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID == 2 {
		return false, nil
	}
	return true, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("room not found with this id")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}

func (m *testDBRepo) UpdateUser(user models.User) error {
	return nil
}

// Authenticate implements [repository.DatabaseRepo].
func (m *testDBRepo) Authenticate(email string, password string) (int, string, error) {
	panic("unimplemented")
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	return []models.Reservation{}, nil
}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	return []models.Reservation{}, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	return models.Reservation{}, nil
}

func (m *testDBRepo) UpdateReservation(reservation models.Reservation) error {
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}
