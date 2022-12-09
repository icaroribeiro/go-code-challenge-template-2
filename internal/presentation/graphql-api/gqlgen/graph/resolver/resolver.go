package resolver

import (
	authservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/auth"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/healthcheck"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/user"
)

type Resolver struct {
	HealthCheckService healthcheckservice.IService
	AuthService        authservice.IService
	UserService        userservice.IService
}

func New(healthCheckService healthcheckservice.IService,
	authService authservice.IService, userService userservice.IService) *Resolver {
	return &Resolver{
		HealthCheckService: healthCheckService,
		AuthService:        authService,
		UserService:        userService,
	}
}
