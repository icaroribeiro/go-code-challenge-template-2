package entity

import (
	"github.com/bluele/factory-go/factory"
	uuid "github.com/satori/go.uuid"
)

// AuthFactory is the function that returns an instance of the auth's domain entity for performing tests.
func AuthFactory(args map[string]interface{}) Auth {
	authFactory := factory.NewFactory(
		Auth{},
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
	})

	return authFactory.MustCreate().(Auth)
}
