// Package repository handles the database dependencies
package repository

import (
	"time"

	"booking-app/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) (int, error)
	SearchAvailabilityByDates(start, end time.Time, roomID int) (int, error)
}
