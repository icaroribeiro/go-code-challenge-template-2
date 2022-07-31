package entity

import (
	"time"

	"github.com/bluele/factory-go/factory"
	fake "github.com/brianvoe/gofakeit/v5"
	datastoreentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/entity"
	uuid "github.com/satori/go.uuid"
)

// NewUser is the function that returns an instance of the the user's datastore entity for performing tests.
func NewUser(args map[string]interface{}) datastoreentity.User {
	userFactory := factory.NewFactory(
		datastoreentity.User{},
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

	return userFactory.MustCreate().(datastoreentity.User)
}
