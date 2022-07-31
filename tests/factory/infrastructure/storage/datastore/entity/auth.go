package entity

import (
	"time"

	"github.com/bluele/factory-go/factory"
	datastoreentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/entity"
	uuid "github.com/satori/go.uuid"
)

// NewAuth is the function that returns an instance of the auth's datastore entity for performing tests.
func NewAuth(args map[string]interface{}) datastoreentity.Auth {
	authFactory := factory.NewFactory(
		datastoreentity.Auth{},
	).Attr("ID", func(fArgs factory.Args) (interface{}, error) {
		id := uuid.NewV4()

		if val, ok := args["id"]; ok {
			id = val.(uuid.UUID)
		}

		return id, nil
	}).Attr("UserID", func(fArgs factory.Args) (interface{}, error) {
		userID := uuid.NewV4()

		if val, ok := args["userID"]; ok {
			userID = val.(uuid.UUID)
		}

		return userID, nil
	}).Attr("CreatedAt", func(fArgs factory.Args) (interface{}, error) {
		createdAt := time.Now()

		if val, ok := args["createdAt"]; ok {
			createdAt = val.(time.Time)
		}

		return createdAt, nil
	})

	return authFactory.MustCreate().(datastoreentity.Auth)
}
