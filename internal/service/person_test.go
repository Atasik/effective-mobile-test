package service

import (
	"encoding/json"
	"errors"
	"fio/internal/domain"
	mock_repository "fio/internal/repository/mocks"
	mock_cache "fio/pkg/cache/mocks"
	mock_profiler "fio/pkg/profiler/mocks"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPersonService_GetAll(t *testing.T) {
	type mockBehaviour func(rp *mock_repository.MockPersonRepo, c *mock_cache.MockCache, t time.Duration, opts domain.PersonsQuery)

	tests := []struct {
		name          string
		inputOpts     domain.PersonsQuery
		mockBehaviour mockBehaviour
		want          []domain.Person
		wantErr       bool
	}{
		{
			name:      "DB OK",
			inputOpts: domain.PersonsQuery{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, c *mock_cache.MockCache, t time.Duration, opts domain.PersonsQuery) {
				c.EXPECT().Get(fmt.Sprintf("getPersons:%v", opts)).Return([]byte{}, errors.New("something went wrong"))
				rp.EXPECT().GetAll(opts).Return([]domain.Person{{ID: 1, Name: "Test", Surname: "Test", Age: 22, Gender: "male", Nationality: "RU"}}, nil)
				personBytes, _ := json.Marshal([]domain.Person{{ID: 1, Name: "Test", Surname: "Test", Age: 22, Gender: "male", Nationality: "RU"}}) //nolint:errcheck
				c.EXPECT().Set(fmt.Sprintf("getPersons:%v", opts), personBytes, t).Return(nil)
			},
			want: []domain.Person{{ID: 1, Name: "Test", Surname: "Test", Age: 22, Gender: "male", Nationality: "RU"}},
		},
		{
			name:      "CacheOK",
			inputOpts: domain.PersonsQuery{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, c *mock_cache.MockCache, t time.Duration, opts domain.PersonsQuery) {
				personBytes, _ := json.Marshal([]domain.Person{{ID: 1, Name: "Test", Surname: "Test", Age: 22, Gender: "male", Nationality: "RU"}})
				c.EXPECT().Get(fmt.Sprintf("getPersons:%v", opts)).Return(personBytes, nil)
			},
			want: []domain.Person{{ID: 1, Name: "Test", Surname: "Test", Age: 22, Gender: "male", Nationality: "RU"}},
		},
		{
			name:      "Cache Error",
			inputOpts: domain.PersonsQuery{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, c *mock_cache.MockCache, t time.Duration, opts domain.PersonsQuery) {
				c.EXPECT().Get(fmt.Sprintf("getPersons:%v", opts)).Return([]byte{1}, nil)
			},
			wantErr: true,
		},
		{
			name:      "DB Error",
			inputOpts: domain.PersonsQuery{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, c *mock_cache.MockCache, t time.Duration, opts domain.PersonsQuery) {
				c.EXPECT().Get(fmt.Sprintf("getPersons:%v", opts)).Return([]byte{}, errors.New("something went wrong"))
				rp.EXPECT().GetAll(opts).Return(nil, errors.New("something went wrong"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		c := gomock.NewController(t)
		defer c.Finish()

		repoPerson := mock_repository.NewMockPersonRepo(c)
		cache := mock_cache.NewMockCache(c)
		cacheTTL := 30 * time.Second //nolint:gomnd
		test.mockBehaviour(repoPerson, cache, cacheTTL, test.inputOpts)

		personService := NewPersonService(repoPerson, cache, nil, cacheTTL)

		got, err := personService.GetAll(test.inputOpts)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		}
	}
}

func TestPersonService_Add(t *testing.T) {
	type mockBehaviour func(rp *mock_repository.MockPersonRepo, np *mock_profiler.MockProfiler, person domain.Person)

	tests := []struct {
		name          string
		inputPerson   domain.Person
		mockBehaviour mockBehaviour
		want          int
		wantErr       bool
	}{
		{
			name:        "OK",
			inputPerson: domain.Person{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, np *mock_profiler.MockProfiler, person domain.Person) {
				np.EXPECT().AgifyPerson(person.Name).Return(6106, nil)
				np.EXPECT().GenderizePerson(person.Name).Return("blop", nil)
				np.EXPECT().NationalizePerson(person.Name).Return("blop", nil)
				person.Age = 6106
				person.Gender = "blop"
				person.Nationality = "blop"
				rp.EXPECT().Add(person).Return(1, nil)
			},
			want: 1,
		},
		{
			name:        "DB Error",
			inputPerson: domain.Person{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, np *mock_profiler.MockProfiler, person domain.Person) {
				np.EXPECT().AgifyPerson(person.Name).Return(6106, nil)
				np.EXPECT().GenderizePerson(person.Name).Return("blop", nil)
				np.EXPECT().NationalizePerson(person.Name).Return("blop", nil)
				person.Age = 6106
				person.Gender = "blop"
				person.Nationality = "blop"
				rp.EXPECT().Add(person).Return(0, errors.New("something went wrong"))
			},
			wantErr: true,
		},
		{
			name:        "Agify Error",
			inputPerson: domain.Person{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, np *mock_profiler.MockProfiler, person domain.Person) {
				np.EXPECT().AgifyPerson(person.Name).Return(0, errors.New("something went wrong"))
			},
			wantErr: true,
		},
		{
			name:        "Genderize Error",
			inputPerson: domain.Person{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, np *mock_profiler.MockProfiler, person domain.Person) {
				np.EXPECT().AgifyPerson(person.Name).Return(6106, nil)
				np.EXPECT().GenderizePerson(person.Name).Return("", errors.New("something went wrong"))
			},
			wantErr: true,
		},
		{
			name:        "Nationalize Error",
			inputPerson: domain.Person{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, np *mock_profiler.MockProfiler, person domain.Person) {
				np.EXPECT().AgifyPerson(person.Name).Return(6106, nil)
				np.EXPECT().GenderizePerson(person.Name).Return("blop", nil)
				np.EXPECT().NationalizePerson(person.Name).Return("", errors.New("something went wrong"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		c := gomock.NewController(t)
		defer c.Finish()

		repoPerson := mock_repository.NewMockPersonRepo(c)
		nameProfiler := mock_profiler.NewMockProfiler(c)
		test.mockBehaviour(repoPerson, nameProfiler, test.inputPerson)

		personService := NewPersonService(repoPerson, nil, nameProfiler, 0) //nolint:gomnd

		got, err := personService.Add(test.inputPerson)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		}
	}
}

func TestPersonService_Delete(t *testing.T) {
	type mockBehaviour func(rp *mock_repository.MockPersonRepo, personID int)

	tests := []struct {
		name          string
		inputID       int
		mockBehaviour mockBehaviour
		want          bool
		wantErr       bool
	}{
		{
			name:    "OK",
			inputID: 1,
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, personID int) {
				rp.EXPECT().Delete(personID).Return(true, nil)
			},
			want: true,
		},
		{
			name:    "Not deleted",
			inputID: 1,
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, personID int) {
				rp.EXPECT().Delete(personID).Return(false, nil)
			},
			want: false,
		},
		{
			name:    "DB Error",
			inputID: 1,
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, personID int) {
				rp.EXPECT().Delete(personID).Return(false, errors.New("something went wrong"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		c := gomock.NewController(t)
		defer c.Finish()

		repoPerson := mock_repository.NewMockPersonRepo(c)
		test.mockBehaviour(repoPerson, test.inputID)

		personService := NewPersonService(repoPerson, nil, nil, 0) //nolint:gomnd

		got, err := personService.Delete(test.inputID)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		}
	}
}

func TestPersonService_Update(t *testing.T) {
	type mockBehaviour func(rp *mock_repository.MockPersonRepo, personID int, updatePerson domain.UpdatePersonInput)

	tests := []struct {
		name             string
		inputID          int
		inputUpdateInput domain.UpdatePersonInput
		mockBehaviour    mockBehaviour
		want             bool
		wantErr          bool
	}{
		{
			name:             "OK",
			inputID:          1,
			inputUpdateInput: domain.UpdatePersonInput{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, personID int, updatePerson domain.UpdatePersonInput) {
				rp.EXPECT().Update(personID, updatePerson).Return(true, nil)
			},
			want: true,
		},
		{
			name:             "Not updated",
			inputID:          1,
			inputUpdateInput: domain.UpdatePersonInput{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, personID int, updatePerson domain.UpdatePersonInput) {
				rp.EXPECT().Update(personID, updatePerson).Return(false, nil)
			},
			want: false,
		},
		{
			name:             "DB Error",
			inputID:          1,
			inputUpdateInput: domain.UpdatePersonInput{},
			mockBehaviour: func(rp *mock_repository.MockPersonRepo, personID int, updatePerson domain.UpdatePersonInput) {
				rp.EXPECT().Update(personID, updatePerson).Return(false, errors.New("something went wrong"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		c := gomock.NewController(t)
		defer c.Finish()

		repoPerson := mock_repository.NewMockPersonRepo(c)
		test.mockBehaviour(repoPerson, test.inputID, test.inputUpdateInput)

		personService := NewPersonService(repoPerson, nil, nil, 0) //nolint:gomnd

		got, err := personService.Update(test.inputID, test.inputUpdateInput)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		}
	}
}
