package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidOffset = errors.New("invalid offset")
	ErrInvalidLimit  = errors.New("invalid limit")
	ErrNotFound      = errors.New("not found")
)

func NewPostgresqlDB(host, port, username, dbname, password, sslmode string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, username, dbname, password, sslmode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ParsePostgresError(err error) error {
	if err == nil {
		return nil
	}

	pgErr, ok := err.(*pq.Error)
	if ok {
		if pgErr.Code == "23505" {
			return ErrAlreadyExists
		}
		if pgErr.Code == "2201X" {
			return ErrInvalidOffset
		}
		if pgErr.Code == "2201W" {
			return ErrInvalidLimit
		}
	}

	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	return err
}
