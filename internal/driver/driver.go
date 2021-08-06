package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//DB holds database pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxDbConnections = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

//ConnectSQL creates db pool to connect to Postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetConnMaxIdleTime(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)
	d.SetMaxOpenConns(maxDbConnections)

	dbConn.SQL = d

	err = TestDb(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

//TestDb tests databse connection
func TestDb(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

//NewDatabase creates ne DB for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
