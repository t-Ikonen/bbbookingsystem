package dbrepo

import (
	"errors"
	"time"

	"github.com/t-Ikonen/bbbookingsystem/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts reservation to DB
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

//InsertRoomRestriction inserts room restriction into db (restriction == reservations)
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	return nil
}

//SearchAvailabilityByDatesByRoomId return false if given room has no availability, return true if availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room
	return rooms, nil
}

// GetRoomNameById gets room by ID
func (m *testDBRepo) GetRoomNameById(id int) (models.Room, error) {
	var room models.Room
	// 3 is number of rooms in DB for renting
	if id > 3 {
		return room, errors.New("some error")
	}
	return room, nil
}

// GetUsedById(id int) (models.User, error)
// UpdateUser(user models.User) error
// Authenticate(email, testPassword string) (int, string, error)
func (m *testDBRepo) GetUsedById(id int) (models.User, error) {
	var user models.User
	return user, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {

	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

//AllReservations gets a slice  of  all the reservations for admin use
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservation []models.Reservation
	return reservation, nil
}

//AllNewReservations gets a slice  of  all the unprocessed reservations for admin use
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservation []models.Reservation
	return reservation, nil
}

//GetReservationById get one unprocessed reservation by reservation ID
func (m *testDBRepo) GetReservationById(id int) (models.Reservation, error) {
	var res models.Reservation
	return res, nil

}

//UpdateReservation updates reservation
func (m *testDBRepo) UpdateReservation(u models.Reservation) error {

	return nil
}

//deletes one reservation by ID
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}
func (m *testDBRepo) UpdatePrcessed(id, processed int) error {

	return nil
}

//return all rooms
func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

//return restriction for a room by ID
func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction

	return restrictions, nil
}
