// Package dbrepo handle the database operations using postgres
package dbrepo

import (
	"context"
	"time"

	"booking-app/internal/models"
)

// AllUsers implements [repository.DatabaseRepo].
func (m *postgresDBRepo) AllUsers() bool {
	panic("unimplemented")
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name,last_name,email,phone,start_date,end_date,room_id,created_at,updated_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	err := m.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, res.CreatedAt, res.UpdatedAt).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction record into the database
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var id int
	stmt := `insert into room_restrictions(start_date,end_date,room_id,reservation_id,restriction_id,created_at,updated_at) values($1,$2,$3,$4,$5,$6,$7) returning id`
	err := m.DB.QueryRowContext(ctx, stmt, res.StartDate, res.EndDate, res.RoomID, res.ReservationID, res.RestrictionID, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `select count(id) from room_restrictions where room_id=$1 and $2 < end_date and $3 > start_date;`
	var numOfRows int

	row := m.DB.QueryRowContext(ctx, query, start, end, roomID)
	err := row.Scan(&numOfRows)
	if err != nil {
		return false, err
	}

	if numOfRows == 0 {
		return false, nil
	}

	return true, nil
}

func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `
	select r.id,r.room_name from rooms r
	where r.id not in (select rr.room_id from room_restrictions where $1<rr.end_date and $2 > rr.start_date)
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
