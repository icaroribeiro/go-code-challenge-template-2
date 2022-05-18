package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/graph/model"
)

func (r *queryResolver) GetHealthCheck(ctx context.Context) (*model.HealthCheck, error) {
	healthcheck := model.HealthCheck{}

	if err := r.HealthCheckService.GetStatus(); err != nil {
		return nil, err
	}

	healthcheck.Status = "everything is up and running"
	return &healthcheck, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
