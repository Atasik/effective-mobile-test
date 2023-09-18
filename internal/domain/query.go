package domain

import "context"

type PaginationKey string

const PaginationContextKey PaginationKey = "query_pagination"

type PaginationQuery struct {
	Limit  int
	Offset int
}

type PersonFiltersQuery struct {
	Name        *string `schema:"name" example:"Vladimir"`
	Surname     *string `schema:"surname" example:"Davydov"`
	Patronymic  *string `schema:"patronymic" example:"Viktorovych"`
	Age         *int    `schema:"age" example:"35"`
	Gender      *string `schema:"gender" example:"male"`
	Nationality *string `schema:"nationality" example:"RU"`
}

type PersonsQuery struct {
	PaginationQuery
	PersonFiltersQuery
}

func PaginationFromContext(ctx context.Context) (*PaginationQuery, error) {
	options, ok := ctx.Value(PaginationContextKey).(*PaginationQuery)
	if !ok || options == nil {
		return nil, ErrNoOptions
	}
	return options, nil
}

func ContextWithPagination(ctx context.Context, opts *PaginationQuery) context.Context {
	return context.WithValue(ctx, PaginationContextKey, opts)
}
