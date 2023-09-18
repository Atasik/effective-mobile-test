package domain

import (
	"errors"
)

var (
	ErrNoOptions      = errors.New("no options")
	ErrUpdateNoFields = errors.New("update structure has no values")
)

type Person struct {
	ID          int     `db:"id" example:"1"`
	Name        string  `json:"name" db:"name" schema:"name" validate:"required" example:"Dmitriy"`
	Surname     string  `json:"surname" db:"surname" schema:"surname" validate:"required" example:"Ushakov"`
	Patronymic  *string `json:"patronymic,omitempty" db:"patronymic" schema:"patronymic" example:"Vasilevich"`
	Age         int     `json:"age" db:"age" schema:"age" example:"22"`
	Gender      string  `json:"gender" db:"gender" schema:"gender" example:"male"`
	Nationality string  `json:"nationality" db:"nationality" schema:"nationality"`
}

type UpdatePersonInput struct {
	Name        *string `json:"name" db:"name" example:"Alexey"`
	Surname     *string `json:"surname" db:"surname" example:"Yakovlev"`
	Patronymic  *string `json:"patronymic" db:"patronymic" example:"Vladimirovich"`
	Age         *int    `json:"age" db:"age" example:"22"`
	Gender      *string `json:"gender" db:"gender" example:"male"`
	Nationality *string `json:"nationality" db:"nationality" example:"RU"`
}

func (i UpdatePersonInput) Validate() error {
	if i.Name == nil && i.Surname == nil && i.Patronymic == nil && i.Age == nil && i.Gender == nil && i.Nationality == nil {
		return ErrUpdateNoFields
	}

	return nil
}
