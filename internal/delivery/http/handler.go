package v1

import (
	"fio/internal/delivery/http/graph"
	"net/http"

	"fio/internal/service"

	_ "fio/docs"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// На гит-хаб странице пакета gorilla/schema написано, что можно объвлять,
// как глобальную переменную, потому что декодер кэширует мета-данные струтур :)
var decoder = schema.NewDecoder() //nolint:gochecknoglobals

const (
	appJSON = "application/json"
)

type Handler struct {
	services  *service.Service
	validator *validator.Validate
	logger    *zap.SugaredLogger
}

func NewHandler(services *service.Service, validator *validator.Validate, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		services:  services,
		validator: validator,
		logger:    logger,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	decoder.IgnoreUnknownKeys(true)
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/api/persons", h.paginationMiddleware(h.getPersons)).Methods("GET")
	r.HandleFunc("/api/person", h.addPerson).Methods("POST")
	r.HandleFunc("/api/person/{personID}", h.deletePerson).Methods("DELETE")
	r.HandleFunc("/api/person/{personID}", h.updatePerson).Methods("PUT")

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(h.services, h.validator, h.logger)}))

	r.Handle("/", playground.Handler("Test Effective-Mobile", "/query"))
	r.Handle("/query", srv)

	mux := h.accessLogMiddleware(r)
	mux = h.panicMiddleware(mux)

	return mux
}
