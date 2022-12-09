package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/entity"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql-api/gqlgen/graph/generated"
)

func (r *queryResolver) GetHealthCheck(ctx context.Context) (*entity.HealthCheck, error) {
	if err := r.HealthCheckService.GetStatus(); err != nil {
		return nil, err
	}

	return &entity.HealthCheck{
		Status: "everything is up and running",
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
