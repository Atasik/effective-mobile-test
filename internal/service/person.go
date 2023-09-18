package service

import (
	"encoding/json"
	"fio/internal/domain"
	"fio/internal/repository"
	"fio/pkg/cache"
	"fio/pkg/profiler"
	"fmt"
	"time"
)

type PersonService struct {
	personRepo   repository.PersonRepo
	nameProfiler profiler.Profiler
	cache        cache.Cache

	cacheTTL time.Duration
}

func NewPersonService(personRepo repository.PersonRepo, cache cache.Cache,
	nameProfiler profiler.Profiler, cacheTTL time.Duration) *PersonService {
	return &PersonService{personRepo: personRepo, cache: cache,
		nameProfiler: nameProfiler, cacheTTL: cacheTTL}
}

func (s *PersonService) GetPersons(opts domain.PersonsQuery) ([]domain.Person, error) {
	var persons []domain.Person
	redisKey := fmt.Sprintf("getPersons:%v", opts)
	if value, err := s.cache.Get(redisKey); err == nil {
		if err = json.Unmarshal(value, &persons); err != nil {
			return []domain.Person{}, err
		}
		return persons, nil
	}

	persons, err := s.personRepo.GetAll(opts)
	if err != nil {
		return []domain.Person{}, err
	}

	personsBytes, err := json.Marshal(persons)
	if err != nil {
		return []domain.Person{}, err
	}

	err = s.cache.Set(redisKey, personsBytes, s.cacheTTL)
	return persons, err
}

func (s *PersonService) AddPerson(person domain.Person) (int, error) {
	age, err := s.nameProfiler.AgifyPerson(person.Name)
	if err != nil {
		return 0, err
	}

	gender, err := s.nameProfiler.GenderizePerson(person.Name)
	if err != nil {
		return 0, err
	}

	nationality, err := s.nameProfiler.NationalizePerson(person.Name)
	if err != nil {
		return 0, err
	}

	person.Age = age
	person.Gender = gender
	person.Nationality = nationality
	return s.personRepo.Add(person)
}

func (s *PersonService) DeletePerson(personID int) (bool, error) {
	rowsUpdated, err := s.personRepo.Delete(personID)
	if err != nil {
		return false, err
	}

	return rowsUpdated > 0, nil
}

func (s *PersonService) UpdatePerson(personID int, input domain.UpdatePersonInput) (bool, error) {
	rowsUpdated, err := s.personRepo.Update(personID, input)
	if err != nil {
		return false, err
	}

	return rowsUpdated > 0, nil
}
