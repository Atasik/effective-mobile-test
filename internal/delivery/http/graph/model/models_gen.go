// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewPerson struct {
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
}
