package domain

import "context"

type PaginationKey string

const PaginationContextKey PaginationKey = "query_pagination"

type PaginationQuery struct {
	Limit  int
	Offset int
}

type PersonFiltersQuery struct {
	Name        *string `schema:"name"`
	Surname     *string `schema:"surname"`
	Patronymic  *string `schema:"patronymic"`
	Age         *int    `schema:"age"`
	Gender      *string `schema:"gender"`
	Nationality *string `schema:"nationality"`
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
