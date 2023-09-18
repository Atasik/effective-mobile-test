package repository

import (
	"fio/internal/domain"
	"fio/pkg/database/postgres"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PersonPostgresqlRepository struct {
	db *sqlx.DB
}

func NewPersonPostgresqlRepository(db *sqlx.DB) *PersonPostgresqlRepository {
	return &PersonPostgresqlRepository{db: db}
}

func (repo *PersonPostgresqlRepository) GetAll(opts domain.PersonsQuery) ([]domain.Person, error) {
	var persons []domain.Person
	conValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if opts.PersonFiltersQuery.Age != nil {
		conValues = append(conValues, fmt.Sprintf("age=$%d", argID))
		args = append(args, *opts.PersonFiltersQuery.Age)
		argID++
	}

	if opts.PersonFiltersQuery.Gender != nil {
		conValues = append(conValues, fmt.Sprintf("gender=$%d", argID))
		args = append(args, *opts.PersonFiltersQuery.Gender)
		argID++
	}

	if opts.PersonFiltersQuery.Name != nil {
		conValues = append(conValues, fmt.Sprintf("name=$%d", argID))
		args = append(args, *opts.PersonFiltersQuery.Name)
		argID++
	}

	if opts.PersonFiltersQuery.Nationality != nil {
		conValues = append(conValues, fmt.Sprintf("nationality=$%d", argID))
		args = append(args, *opts.PersonFiltersQuery.Nationality)
		argID++
	}

	if opts.PersonFiltersQuery.Patronymic != nil {
		conValues = append(conValues, fmt.Sprintf("patronymic=$%d", argID))
		args = append(args, *opts.PersonFiltersQuery.Patronymic)
		argID++
	}

	if opts.PersonFiltersQuery.Surname != nil {
		conValues = append(conValues, fmt.Sprintf("surname=$%d", argID))
		args = append(args, *opts.PersonFiltersQuery.Surname)
		argID++
	}
	args = append(args, opts.PaginationQuery.Limit, opts.PaginationQuery.Offset)

	conQuery := strings.Join(conValues, " AND ")

	query := fmt.Sprintf("SELECT * FROM %s LIMIT $%d OFFSET $%d", personsTable, argID, argID+1)
	if len(conValues) > 0 {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT $%d OFFSET $%d", personsTable, conQuery, argID, argID+1)
	}

	if err := repo.db.Select(&persons, query, args...); err != nil {
		return []domain.Person{}, postgres.ParsePostgresError(err)
	}

	return persons, nil
}

func (repo *PersonPostgresqlRepository) Add(person domain.Person) (int, error) {
	var personID int

	query := fmt.Sprintf("INSERT INTO %s (name, surname, patronymic, age, gender, nationality) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id", personsTable)

	row := repo.db.QueryRow(query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
	err := row.Scan(&personID)
	if err != nil {
		return 0, postgres.ParsePostgresError(err)
	}
	return personID, nil
}

func (repo *PersonPostgresqlRepository) Delete(personID int) (bool, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", personsTable)

	res, err := repo.db.Exec(query, personID)
	if err != nil {
		return false, postgres.ParsePostgresError(err)
	}
	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

func (repo *PersonPostgresqlRepository) Update(personID int, input domain.UpdatePersonInput) (bool, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Age != nil {
		setValues = append(setValues, fmt.Sprintf("age=$%d", argID))
		args = append(args, *input.Age)
		argID++
	}

	if input.Gender != nil {
		setValues = append(setValues, fmt.Sprintf("gender=$%d", argID))
		args = append(args, *input.Gender)
		argID++
	}

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argID))
		args = append(args, *input.Name)
		argID++
	}

	if input.Nationality != nil {
		setValues = append(setValues, fmt.Sprintf("nationality=$%d", argID))
		args = append(args, *input.Nationality)
		argID++
	}

	if input.Patronymic != nil {
		setValues = append(setValues, fmt.Sprintf("patronymic=$%d", argID))
		args = append(args, *input.Patronymic)
		argID++
	}

	if input.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argID))
		args = append(args, *input.Surname)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, personsTable, setQuery, argID)
	args = append(args, personID)

	res, err := repo.db.Exec(query, args...)
	if err != nil {
		return false, postgres.ParsePostgresError(err)
	}
	updated, err := res.RowsAffected()
	return updated > 0, err
}
