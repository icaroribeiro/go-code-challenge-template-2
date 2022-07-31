package entity

import (
	"github.com/bluele/factory-go/factory"
	fake "github.com/brianvoe/gofakeit/v5"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
	uuid "github.com/satori/go.uuid"
)

// NewUser is the function that returns an instance of the user's domain model for performing tests.
func NewUser(args map[string]interface{}) domainentity.User {
	userFactory := factory.NewFactory(
		domainentity.User{},
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

	return userFactory.MustCreate().(domainentity.User)
}
