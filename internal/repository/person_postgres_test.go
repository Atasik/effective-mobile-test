package repository

import (
	"fio/internal/domain"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestPersonPostgres_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	sqlxDb := sqlx.NewDb(db, "sqlmock")
	r := NewPersonPostgresqlRepository(sqlxDb)

	tests := []struct {
		name    string
		mock    func()
		input   domain.Person
		want    int
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", personsTable)).
					WithArgs("TEST", "TEST", "TEST", 54, "TEST", "TEST").WillReturnRows(rows)
			},
			input: domain.Person{
				Name:        "TEST",
				Surname:     "TEST",
				Patronymic:  stringPointer("TEST"),
				Age:         54,
				Gender:      "TEST",
				Nationality: "TEST",
			},
			want: 1,
		},
		{
			name: "Empty Field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", personsTable)).
					WithArgs("", "TEST", "TEST", 54, "TEST", "TEST").WillReturnRows(rows)
			},
			input: domain.Person{
				Name:        "",
				Surname:     "TEST",
				Patronymic:  stringPointer("TEST"),
				Age:         54,
				Gender:      "TEST",
				Nationality: "TEST",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Add(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPersonPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	sqlxDb := sqlx.NewDb(db, "sqlmock")
	r := NewPersonPostgresqlRepository(sqlxDb)

	tests := []struct {
		name    string
		mock    func()
		want    bool
		input   int
		wantErr bool
	}{
		{
			name: "Deleted",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("DELETE FROM %s WHERE (.+)", personsTable)).
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			want:  true,
			input: 1,
		},
		{
			name: "Not Deleted",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("DELETE FROM %s WHERE (.+)", personsTable)).
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			input: 1,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Delete(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPersonPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDb := sqlx.NewDb(db, "sqlmock")
	r := NewPersonPostgresqlRepository(sqlxDb)

	type args struct {
		personID int
		input    domain.UpdatePersonInput
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
		want    bool
	}{
		{
			name: "OK_AllFields",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+)", personsTable)).
					WithArgs(25, "new gender", "new name", "new nationality", "new patronymic", "new surname", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				personID: 1,
				input: domain.UpdatePersonInput{
					Age:         intPointer(25),
					Gender:      stringPointer("new gender"),
					Name:        stringPointer("new name"),
					Nationality: stringPointer("new nationality"),
					Patronymic:  stringPointer("new patronymic"),
					Surname:     stringPointer("new surname"),
				},
			},
			want: true,
		},
		{
			name: "OK_NotUpdated",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+)", personsTable)).
					WithArgs(25, "new gender", "new name", "new nationality", "new patronymic", "new surname", 1).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			input: args{
				personID: 1,
				input: domain.UpdatePersonInput{
					Age:         intPointer(25),
					Gender:      stringPointer("new gender"),
					Name:        stringPointer("new name"),
					Nationality: stringPointer("new nationality"),
					Patronymic:  stringPointer("new patronymic"),
					Surname:     stringPointer("new surname"),
				},
			},
			want: false,
		},
		{
			name: "OK_WithoutSurname",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+)", personsTable)).
					WithArgs(25, "new gender", "new name", "new nationality", "new patronymic", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				personID: 1,
				input: domain.UpdatePersonInput{
					Age:         intPointer(25),
					Gender:      stringPointer("new gender"),
					Name:        stringPointer("new name"),
					Nationality: stringPointer("new nationality"),
					Patronymic:  stringPointer("new patronymic"),
				},
			},
			want: true,
		},
		{
			name: "OK_WithoutSurnameAndPatronymic",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+)", personsTable)).
					WithArgs(25, "new gender", "new name", "new nationality", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				personID: 1,
				input: domain.UpdatePersonInput{
					Age:         intPointer(25),
					Gender:      stringPointer("new gender"),
					Name:        stringPointer("new name"),
					Nationality: stringPointer("new nationality"),
				},
			},
			want: true,
		},
		{
			name: "OK_WithoutSurnameAndPatronymicAndNationality",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+)", personsTable)).
					WithArgs(25, "new gender", "new name", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				personID: 1,
				input: domain.UpdatePersonInput{
					Age:    intPointer(25),
					Gender: stringPointer("new gender"),
					Name:   stringPointer("new name"),
				},
			},
			want: true,
		},
		{
			name: "OK_WithoutSurnameAndPatronymicAndNationalityAndName",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+)", personsTable)).
					WithArgs(25, "new gender", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				personID: 1,
				input: domain.UpdatePersonInput{
					Age:    intPointer(25),
					Gender: stringPointer("new gender"),
				},
			},
			want: true,
		},
		{
			name: "OK_WithoutSurnameAndPatronymicAndNationalityAndNameAndGender",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+)", personsTable)).
					WithArgs(25, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				personID: 1,
				input: domain.UpdatePersonInput{
					Age: intPointer(25),
				},
			},
			want: true,
		},
		{
			name: "OK_NoInputFields",
			mock: func() {
				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET WHERE (.+)", personsTable)).
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			input: args{
				personID: 1,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Update(tt.input.personID, tt.input.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func intPointer(i int) *int {
	return &i
}
