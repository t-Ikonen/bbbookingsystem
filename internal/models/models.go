package models

import "time"

//Reservation holds reservation data
// type Reservation struct {
// 	FirstName string
// 	LastName  string
// 	Phone     string
// 	Email     string
// }

//User holds DB model of user table
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

//Room hold room model
type Room struct {
	ID         int
	RoomName   string
	Minibar    bool
	Shower     bool
	PricingId  int
	CreatedAt  time.Time
	ModifiedAt time.Time
}

//Restriction is restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	ModifiedAt      time.Time
}

//Reservations is reservations model
type Reservation struct {
	ID         int
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	StartDate  time.Time
	EndDate    time.Time
	RoomId     int
	CreatedAt  time.Time
	ModifiedAt time.Time
	Room       Room
}

//RoomRestrictions is RoomRestrictions model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndSate       time.Time
	RoomId        int
	ReservationId int
	CreatedAt     time.Time
	ModifiedAt    time.Time
	RestrictionId int
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}
