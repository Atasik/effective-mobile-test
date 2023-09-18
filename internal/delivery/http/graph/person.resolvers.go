package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.37

import (
	"context"
	"fio/internal/delivery/http/graph/model"
	"fio/internal/domain"
)

// AddPerson is the resolver for the addPerson field.
func (r *mutationResolver) AddPerson(ctx context.Context, input model.NewPerson) (int, error) {
	person := domain.Person{
		Name:       input.Name,
		Surname:    input.Surname,
		Patronymic: input.Patronymic,
	}
	r.validator.Struct(person)
	return r.services.AddPerson(person)
}

// DeletePerson is the resolver for the deletePerson field.
func (r *mutationResolver) DeletePerson(ctx context.Context, id int) (bool, error) {
	return r.services.DeletePerson(id)
}

// UpdatePerson is the resolver for the updatePerson field.
func (r *mutationResolver) UpdatePerson(ctx context.Context, id int, input domain.UpdatePersonInput) (bool, error) {
	err := input.Validate()
	if err != nil {
		return false, err
	}
	return r.services.UpdatePerson(id, input)
}

// GetPersons is the resolver for the getPersons field.
func (r *queryResolver) GetPersons(ctx context.Context, filter *domain.PersonFiltersQuery, limit *int, offset *int) ([]domain.Person, error) {
	if filter == nil {
		filter = &domain.PersonFiltersQuery{}
	}
	opts := domain.PersonsQuery{
		PaginationQuery:    domain.PaginationQuery{Limit: *limit, Offset: *offset},
		PersonFiltersQuery: *filter,
	}
	return r.services.GetPersons(opts)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
