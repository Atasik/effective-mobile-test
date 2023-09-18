package v1

import (
	"bytes"
	"context"
	"errors"
	"fio/internal/domain"
	"fio/internal/service"
	mock_service "fio/internal/service/mocks"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestHandler_getPersons(t *testing.T) {
	type mockBehaviour func(su *mock_service.MockPerson, opts domain.PersonsQuery)

	tests := []struct {
		name                 string
		inputPersonsQuery    domain.PersonsQuery
		mockBehaviour        mockBehaviour
		params               map[string]string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:              "OK",
			inputPersonsQuery: domain.PersonsQuery{},
			mockBehaviour: func(su *mock_service.MockPerson, opts domain.PersonsQuery) {
				su.EXPECT().GetAll(opts).Return([]domain.Person{{ID: 1, Name: "Test", Surname: "Test", Age: 5, Gender: "male", Nationality: "RU"}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"data":[{"ID":1,"name":"Test","surname":"Test","age":5,"gender":"male","nationality":"RU"}]}`,
		},
		{
			name:              "Service Error",
			inputPersonsQuery: domain.PersonsQuery{},
			mockBehaviour: func(su *mock_service.MockPerson, opts domain.PersonsQuery) {
				su.EXPECT().GetAll(opts).Return(nil, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			servicePerson := mock_service.NewMockPerson(c)
			test.mockBehaviour(servicePerson, test.inputPersonsQuery)

			services := &service.Service{Person: servicePerson}

			validate := validator.New()
			logger := zap.NewNop().Sugar()
			h := NewHandler(services, validate, logger)

			r := mux.NewRouter()
			r.HandleFunc("/api/persons", h.getPersons).Methods("GET")

			w := httptest.NewRecorder()
			ctx := context.WithValue(context.TODO(), domain.PaginationContextKey, &test.inputPersonsQuery.PaginationQuery)
			req := httptest.NewRequest("GET", "/api/persons", nil).WithContext(ctx)
			q := req.URL.Query()
			for k, v := range test.params {
				q.Add(k, v)
			}
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_addPerson(t *testing.T) {
	type mockBehaviour func(su *mock_service.MockPerson, person domain.Person)

	tests := []struct {
		name                 string
		inputBody            string
		inputPerson          domain.Person
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"alex", "surname":"test"}`,
			inputPerson: domain.Person{
				Name:    "alex",
				Surname: "test",
			},
			mockBehaviour: func(su *mock_service.MockPerson, person domain.Person) {
				su.EXPECT().Add(person).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"name":""}`,
			inputPerson:          domain.Person{},
			mockBehaviour:        func(su *mock_service.MockPerson, person domain.Person) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"bad input"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"name":"alex", "surname":"test"}`,
			inputPerson: domain.Person{
				Name:    "alex",
				Surname: "test",
			},
			mockBehaviour: func(su *mock_service.MockPerson, person domain.Person) {
				su.EXPECT().Add(person).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			servicePerson := mock_service.NewMockPerson(c)
			test.mockBehaviour(servicePerson, test.inputPerson)

			services := &service.Service{Person: servicePerson}

			validate := validator.New()
			logger := zap.NewNop().Sugar()
			h := NewHandler(services, validate, logger)

			r := mux.NewRouter()
			r.HandleFunc("/api/person", h.addPerson).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/person",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", appJSON)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_deletePerson(t *testing.T) {
	type mockBehaviour func(su *mock_service.MockPerson, personID int)

	tests := []struct {
		name                 string
		paramID              string
		inputID              int
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "OK",
			paramID: "1",
			inputID: 1,
			mockBehaviour: func(su *mock_service.MockPerson, personID int) {
				su.EXPECT().Delete(personID).Return(true, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"done"}`,
		},
		{
			name:                 "Bad ID",
			paramID:              "1d",
			inputID:              1,
			mockBehaviour:        func(su *mock_service.MockPerson, personID int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"bad id"}`,
		},
		{
			name:    "Service Error",
			paramID: "1",
			inputID: 1,
			mockBehaviour: func(su *mock_service.MockPerson, personID int) {
				su.EXPECT().Delete(personID).Return(false, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			servicePerson := mock_service.NewMockPerson(c)
			test.mockBehaviour(servicePerson, test.inputID)

			services := &service.Service{Person: servicePerson}

			validate := validator.New()
			logger := zap.NewNop().Sugar()
			h := NewHandler(services, validate, logger)

			r := mux.NewRouter()
			r.HandleFunc("/api/person/{personID}", h.deletePerson).Methods("DELETE")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/person/%s", test.paramID), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_updatePerson(t *testing.T) {
	type mockBehaviour func(su *mock_service.MockPerson, personID int, updatePerson domain.UpdatePersonInput)

	tests := []struct {
		name                 string
		inputBody            string
		paramID              string
		inputID              int
		inputUpdateInput     domain.UpdatePersonInput
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"test"}`,
			paramID:   "1",
			inputID:   1,
			inputUpdateInput: domain.UpdatePersonInput{
				Name: stringPointer("test"),
			},
			mockBehaviour: func(su *mock_service.MockPerson, personID int, updatePerson domain.UpdatePersonInput) {
				su.EXPECT().Update(personID, updatePerson).Return(true, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"done"}`,
		},
		{
			name:                 "No update values",
			inputBody:            `{}`,
			paramID:              "1",
			inputID:              1,
			inputUpdateInput:     domain.UpdatePersonInput{},
			mockBehaviour:        func(su *mock_service.MockPerson, personID int, updatePerson domain.UpdatePersonInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"update structure has no values"}`,
		},
		{
			name:                 "Wrong ID",
			inputBody:            `{"name":"test"}`,
			paramID:              "1d",
			inputID:              1,
			inputUpdateInput:     domain.UpdatePersonInput{},
			mockBehaviour:        func(su *mock_service.MockPerson, personID int, updatePerson domain.UpdatePersonInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"bad id"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"name":"test"}`,
			paramID:   "1",
			inputID:   1,
			inputUpdateInput: domain.UpdatePersonInput{
				Name: stringPointer("test"),
			},
			mockBehaviour: func(su *mock_service.MockPerson, personID int, updatePerson domain.UpdatePersonInput) {
				su.EXPECT().Update(personID, updatePerson).Return(false, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			servicePerson := mock_service.NewMockPerson(c)
			test.mockBehaviour(servicePerson, test.inputID, test.inputUpdateInput)

			services := &service.Service{Person: servicePerson}

			validate := validator.New()
			logger := zap.NewNop().Sugar()
			h := NewHandler(services, validate, logger)

			r := mux.NewRouter()
			r.HandleFunc("/api/person/{personID}", h.updatePerson).Methods("PUT")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/person/%s", test.paramID),
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", appJSON)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func stringPointer(s string) *string {
	return &s
}
