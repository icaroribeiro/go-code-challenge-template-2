package model

import (
	"time"

	"github.com/bluele/factory-go/factory"
	fake "github.com/brianvoe/gofakeit/v5"
	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/model"
	uuid "github.com/satori/go.uuid"
)

// NewUser is the function that returns an instance of the the user's datastore model for performing tests.
func NewUser(args map[string]interface{}) datastoremodel.User {
	userFactory := factory.NewFactory(
		datastoremodel.User{},
	).Attr("ID", func(fArgs factory.Args) (interface{}, error) {
		id := uuid.NewV4()

		if val, ok := args["id"]; ok {
			id = val.(uuid.UUID)
		}

		return id, nil
	}).Attr("Username", func(fArgs factory.Args) (interface{}, error) {
		username := fake.Username()

		if val, ok := args["username"]; ok {
			username = val.(string)
		}

		return username, nil
	}).Attr("CreatedAt", func(fArgs factory.Args) (interface{}, error) {
		createdAt := time.Now()

		if val, ok := args["createdAt"]; ok {
			createdAt = val.(time.Time)
		}

		return createdAt, nil
	}).Attr("UpdatedAt", func(fArgs factory.Args) (interface{}, error) {
		updatedAt := time.Now()

		if val, ok := args["updatedAt"]; ok {
			updatedAt = val.(time.Time)
		}

		return updatedAt, nil
	})

	return userFactory.MustCreate().(datastoremodel.User)
}
