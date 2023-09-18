package service

import (
	"fio/internal/domain"
	"fio/internal/repository"
	"fio/pkg/cache"
	"fio/pkg/profiler"
	"time"
)

type Person interface {
	GetAll(opts domain.PersonsQuery) ([]domain.Person, error)
	Add(person domain.Person) (int, error)
	Delete(personID int) (bool, error)
	Update(personID int, UpdateInput domain.UpdatePersonInput) (bool, error)
}

type Service struct {
	Person
}

type Dependencies struct {
	Repos        *repository.Repository
	Cache        cache.Cache
	NameProfiler profiler.Profiler
	CacheTTL     time.Duration
}

func NewService(deps Dependencies) *Service {
	return &Service{
		Person: NewPersonService(deps.Repos, deps.Cache, deps.NameProfiler, deps.CacheTTL),
	}
}
