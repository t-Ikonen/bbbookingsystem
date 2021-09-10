package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/t-Ikonen/bbbookingsystem/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts reservation to DB
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	var newId int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `INSERT INTO 
				reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			VALUES
				 ($1,$2,$3,$4,$5,$6,$7,$8,$9) 
			RETURNING id`

	err := m.DB.QueryRowContext(ctx, stmt,
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

//InsertRoomRestriction inserts room restriction into db (restriction == reservations)
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO 
				room_restrictions (start_date, end_date, room_id, reservation_id,  restriction_id, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndSate,
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

//SearchAvailabilityByDatesByRoomId return false if given room has no availability, return true if availability
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		SELECT
			COUNT(id)
		FROM
			room_restrictions 
		WHERE
			$1 < end_date AND
			$2 > start_date AND
			room_id = $3`

	var numRows int

	row := m.DB.QueryRowContext(ctx, stmt, start, end, roomID)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room

	query := `
		SELECT
			r.id, r.room_name 
		FROM
			rooms r 
		WHERE r.id NOT IN
			(SELECT 
				rr.room_id
			FROM
				room_restrictions as rr
			WHERE $1 < rr.end_date AND $2 > rr.start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {

		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomNameById gets room by ID
func (m *postgresDBRepo) GetRoomNameById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		SELECT
			r.id, r.room_name, r.shower, r.minibar, r.pricing_id
		FROM
			rooms AS r
		WHERE 
			r.id= $1 `

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.Shower,
		&room.Minibar,
		&room.PricingId,
	)
	if err != nil {
		return room, err
	}
	return room, nil
}

//GetUsedById returns a user
func (m *postgresDBRepo) GetUsedById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT
			id, first_name, last_name, email, password, access_level
		FROM
			users
		WHERE id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
	)
	if err != err {
		return user, err
	}
	return user, nil
}

//UpdateUser updates used data
func (m *postgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE 
			users
		SET
			first_name = $1, last_name = %2, email = %3, access_level = %4, updated_at = %5
		`

	_, err := m.DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

//Authenticate authenticates user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "SELECT id, password from users WHERE email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password!")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil
}
