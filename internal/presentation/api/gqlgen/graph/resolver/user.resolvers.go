package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/entity"
)

func (r *queryResolver) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	domainUsers, err := r.UserService.WithDBTrx(nil).GetAll()
	if err != nil {
		return nil, err
	}

	users := entity.Users{}
	users.FromDomain(domainUsers)

	allUsers := []*entity.User{}
	for i := range users {
		allUsers = append(allUsers, &users[i])
	}

	return allUsers, nil
}
