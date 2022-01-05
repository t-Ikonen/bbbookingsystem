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
	if err != nil {
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
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil
}

//AllReservations gets a slice  of  all the reservations for admin use
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation []models.Reservation

	query := ` 
		SELECT 
			r.id, r.first_name, r.last_name, r.phone, r.email, r.start_date, r.end_date, r.room_id,
			r.created_at, r.updated_at, r.processed,  rm.id, rm.room_name
		FROM
			reservations as r
		LEFT JOIN
			rooms as rm 
		ON 
			(r.room_id = rm.id) 
		ORDER BY
			r.start_date
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservation, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.ModifiedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservation, err
		}
		reservation = append(reservation, i)
	}
	if err = rows.Err(); err != nil {
		return reservation, err
	}
	return reservation, nil
}

//AllNewReservations gets a slice  of  all the unprocessed reservations for admin use
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation []models.Reservation

	query := ` 
		SELECT 
			r.id, r.first_name, r.last_name, r.phone, r.email, r.start_date, r.end_date, r.room_id,
			r.created_at, r.updated_at, rm.id, rm.room_name
		FROM
			reservations as r
		LEFT JOIN
			rooms as rm 
		ON 
			(r.room_id = rm.id) 
		WHERE
			r.processed=0
		ORDER BY
			r.start_date
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservation, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.ModifiedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservation, err
		}
		reservation = append(reservation, i)
	}
	if err = rows.Err(); err != nil {
		return reservation, err
	}
	return reservation, nil
}

//GetReservationById get one unprocessed reservation by reservation ID
func (m *postgresDBRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var res models.Reservation
	query := `
		SELECT 
			r.id, r.first_name, r.last_name, r.phone, r.email, r.start_date, r.end_date, r.room_id,
			r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
		FROM
			reservations as r
		LEFT JOIN
			rooms as rm 
		ON 
			(r.room_id = rm.id) 
		WHERE
			r.id = $1

		`
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomId,
		&res.CreatedAt,
		&res.ModifiedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

//UpdateReservation updates reservation
func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE 
			reservations
		SET
			first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
		WHERE
			id = $6
		`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

//deletes one reservation by ID
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE FROM
			reservations
		WHERE
			id=$1
	`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

//updates one reservation by ID to be processed
func (m *postgresDBRepo) UpdatePrcessed(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE 
			reservations
		SET
			processed= $1
		WHERE
			id=$2
	`
	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}
	return nil
}

//return all rooms
func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		SELECT * FROM
			rooms
		ORDER BY
			room_name
	`
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return rooms, err
	}

	defer rows.Close()

	for rows.Next() {
		var rm models.Room
		err := rows.Scan(
			&rm.ID,
			&rm.RoomName,
			&rm.Shower,
			&rm.Minibar,
			&rm.PricingId,
			&rm.CreatedAt,
			&rm.ModifiedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, rm)
	}

	return rooms, nil
}

//return restriction for a room by ID and date range
func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `
		SELECT 
			id, coalesce (reservation_id, 0), restriction_id, room_id, start_date, end_date
		FROM
			room_restrictions
		WHERE
			$1 < end_date 
		AND
			$2 >= start_date
		AND
			room_id = $3
	`
	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID,
			&r.ReservationId,
			&r.RestrictionId,
			&r.RoomId,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return nil, err
		}

		restrictions = append(restrictions, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return restrictions, nil
}
