package resolver

import (
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/service/healthcheck"
)

type Resolver struct {
	HealthCheckService healthcheckservice.IService
}

func NewResolver(healthCheckService healthcheckservice.IService) *Resolver {
	return &Resolver{
		HealthCheckService: healthCheckService,
	}
}
