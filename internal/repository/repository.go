package repository

import (
	"fio/internal/domain"

	"github.com/jmoiron/sqlx"
)

const (
	personsTable = "persons"
)

type PersonRepo interface {
	GetAll(opts domain.PersonsQuery) ([]domain.Person, error)
	Add(person domain.Person) (int, error)
	Delete(personID int) (bool, error)
	Update(personID int, input domain.UpdatePersonInput) (bool, error)
}

type Repository struct {
	PersonRepo
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		PersonRepo: NewPersonPostgresqlRepository(db),
	}
}
