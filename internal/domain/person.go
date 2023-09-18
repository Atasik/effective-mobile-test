package domain

import (
	"errors"
)

var (
	ErrNoOptions      = errors.New("no options")
	ErrUpdateNoFields = errors.New("update structure has no values")
)

type Person struct {
	ID          int     `db:"id"`
	Name        string  `json:"name" db:"name" schema:"name" validate:"required"`
	Surname     string  `json:"surname" db:"surname" schema:"surname" validate:"required"`
	Patronymic  *string `json:"patronymic,omitempty" db:"patronymic" schema:"patronymic"`
	Age         int     `json:"age" db:"age" schema:"age"`
	Gender      string  `json:"gender" db:"gender" schema:"gender"`
	Nationality string  `json:"nationality" db:"nationality" schema:"nationality"`
}

type UpdatePersonInput struct {
	Name        *string `json:"name" db:"name"`
	Surname     *string `json:"surname" db:"surname"`
	Patronymic  *string `json:"patronymic" db:"patronymic"`
	Age         *int    `json:"age" db:"age"`
	Gender      *string `json:"gender" db:"gender"`
	Nationality *string `json:"nationality" db:"nationality"`
}

func (i UpdatePersonInput) Validate() error {
	if i.Name == nil && i.Surname == nil && i.Patronymic == nil && i.Age == nil && i.Gender == nil && i.Nationality == nil {
		return ErrUpdateNoFields
	}

	return nil
}
