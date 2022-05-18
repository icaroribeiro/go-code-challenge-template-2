package resolver

import (
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/healthcheck"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/user"
)

type Resolver struct {
	HealthCheckService healthcheckservice.IService
	UserService        userservice.IService
}

func NewResolver(healthCheckService healthcheckservice.IService, userService userservice.IService) *Resolver {
	return &Resolver{
		HealthCheckService: healthCheckService,
		UserService:        userService,
	}
}
