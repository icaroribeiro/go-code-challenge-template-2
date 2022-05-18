package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/model"
)

func (r *queryResolver) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	domainUsers, err := r.UserService.WithDBTrx(nil).GetAll()
	if err != nil {
		return nil, err
	}

	users := model.Users{}
	users.FromDomain(domainUsers)

	allUsers := []*model.User{}
	for i := range users {
		allUsers = append(allUsers, &users[i])
	}

	return allUsers, nil
}
