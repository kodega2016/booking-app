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
