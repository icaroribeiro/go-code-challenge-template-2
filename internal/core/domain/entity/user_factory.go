package entity

import (
	"github.com/bluele/factory-go/factory"
	fake "github.com/brianvoe/gofakeit/v5"
	uuid "github.com/satori/go.uuid"
)

// UserFactory is the function that returns an instance of the user's domain entity for performing tests.
func UserFactory(args map[string]interface{}) User {
	userFactory := factory.NewFactory(
		User{},
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
	})

	return userFactory.MustCreate().(User)
}
