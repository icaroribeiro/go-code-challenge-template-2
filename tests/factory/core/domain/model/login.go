package model

import (
	"github.com/bluele/factory-go/factory"
	fake "github.com/brianvoe/gofakeit/v5"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	uuid "github.com/satori/go.uuid"
)

// NewLogin is the function that returns an instance of the login's domain model for performing tests.
func NewLogin(args map[string]interface{}) domainmodel.Login {
	loginFactory := factory.NewFactory(
		domainmodel.Login{},
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
	}).Attr("Username", func(fArgs factory.Args) (interface{}, error) {
		username := fake.Username()

		if val, ok := args["username"]; ok {
			username = val.(string)
		}

		return username, nil
	}).Attr("Password", func(fArgs factory.Args) (interface{}, error) {
		password := fake.Password(true, true, true, false, false, 8)

		if val, ok := args["password"]; ok {
			password = val.(string)
		}

		return password, nil
	})

	return loginFactory.MustCreate().(domainmodel.Login)
}
