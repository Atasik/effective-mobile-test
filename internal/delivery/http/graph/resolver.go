package graph

import (
	"fio/internal/service"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	services  *service.Service
	validator *validator.Validate
	logger    *zap.SugaredLogger
}

func NewResolver(services *service.Service, validator *validator.Validate, logger *zap.SugaredLogger) *Resolver {
	return &Resolver{
		services:  services,
		validator: validator,
		logger:    logger,
	}
}
