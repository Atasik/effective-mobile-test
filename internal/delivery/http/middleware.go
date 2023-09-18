package v1

import (
	"fio/internal/domain"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultPage  = 1
	defaultLimit = 25
	maxLimit     = 50
)

func (h *Handler) panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Infow("panic middleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				h.logger.Infow("recovered", err)
				newErrorResponse(w, `Internal server error`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Infow("access log middleware", r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		h.logger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Since(start),
		)
	})
}

func (h *Handler) paginationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Infow("query middleware", r.URL.Path)
		limit := parseInt(r.URL.Query().Get("limit"), defaultLimit)
		page := parseInt(r.URL.Query().Get("page"), defaultPage)

		if limit < 1 {
			limit = defaultLimit
		}

		if limit > maxLimit {
			limit = maxLimit
		}

		if page < 1 {
			page = defaultPage
		}

		options := &domain.PaginationQuery{
			Limit:  limit,
			Offset: (page - 1) * limit,
		}
		ctx := domain.ContextWithPagination(r.Context(), options)
		next(w, r.WithContext(ctx))
	}
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}
