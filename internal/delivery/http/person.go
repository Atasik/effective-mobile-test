package v1

import (
	"encoding/json"
	"fio/internal/domain"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// @Summary Get Persons
// @Tags person
// @ID	 get-persons
// @Product json
// @Param   filter  query domain.PersonFiltersQuery false "Query params"
// @Param   limit   query int false "limit" Enums(10, 25, 50)
// @Param   page  query int false "page"
// @Success	200		    {object}	getPersonsResponse
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/persons [get]
func (h *Handler) getPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	var filter domain.PersonFiltersQuery
	err := decoder.Decode(&filter, r.URL.Query())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	pagination, err := domain.PaginationFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	opts := domain.PersonsQuery{PaginationQuery: *pagination, PersonFiltersQuery: filter}
	persons, err := h.services.GetPersons(opts)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetPersonsResponse(w, persons, http.StatusOK)
}

// @Summary Add Person
// @Tags person
// @ID	 add-person
// @Accept json
// @Product json
// @Param   input body domain.Person true "Person content"
// @Success	200		    {integer}	integer     "id"
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/person [post]
func (h *Handler) addPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	if r.Header.Get("Content-Type") != appJSON {
		newErrorResponse(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, "bad input", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var person domain.Person
	err = json.Unmarshal(body, &person)
	if err != nil {
		newErrorResponse(w, "can't unpack payload", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(person)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.services.AddPerson(person)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Infof("Person with ID %d was added", id)

	resp, err := json.Marshal(idResponse{id})
	if err != nil {
		newErrorResponse(w, `can't create payload`, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		newErrorResponse(w, `can't write resp`, http.StatusInternalServerError)
		return
	}
}

// @Summary Delete Person
// @Tags person
// @ID	 delete-person
// @Product json
// @Param		personID	path		integer			true	"ID of person to delete"
// @Success	200		    {object}	statusResponse
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/person/{personID} [delete]
func (h *Handler) deletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	vars := mux.Vars(r)
	personID, err := strconv.Atoi(vars["personID"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	deleted, err := h.services.DeletePerson(personID)
	if err != nil {
		newErrorResponse(w, "DB error", http.StatusInternalServerError)
		return
	}
	if deleted {
		h.logger.Infof("Person with ID %d was deleted: %v", personID, deleted)
	}

	newStatusReponse(w, "done", http.StatusOK)
}

// @Summary Update Person
// @Tags person
// @ID	 update-person
// @Accept json
// @Product json
// @Param		personID	path		integer			true	"ID of person to update"
// @Param		input		body		domain.UpdatePersonInput	false	"Update Person content"
// @Success	200		    {object}	statusResponse
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/person/{personID} [put]
func (h *Handler) updatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	if r.Header.Get("Content-Type") != appJSON {
		newErrorResponse(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, "bad input", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	vars := mux.Vars(r)
	personID, err := strconv.Atoi(vars["personID"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	var inp domain.UpdatePersonInput
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &inp)
	if err != nil {
		newErrorResponse(w, "can't unpack payload", http.StatusBadRequest)
		return
	}
	err = inp.Validate()
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.services.UpdatePerson(personID, inp)
	if err != nil {
		newErrorResponse(w, "DB error", http.StatusInternalServerError)
		return
	}
	if updated {
		h.logger.Infof("Person with ID %d was updated: %v", personID, updated)
	}

	newStatusReponse(w, "done", http.StatusOK)
}
