package perentity

import (
	"time"

	"github.com/bluele/factory-go/factory"
	fake "github.com/brianvoe/gofakeit/v5"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/security"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// LoginFactory is the function that returns an instance of the the login's datastore entity for performing tests.
func LoginFactory(args map[string]interface{}) Login {
	loginFactory := factory.NewFactory(
		Login{},
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

		security := securitypkg.New()

		bytes, err := security.HashPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		hashedPassword := string(bytes)
		if val, ok := args["password"]; ok {
			hashedPassword = val.(string)
		}

		return hashedPassword, nil
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

	return loginFactory.MustCreate().(Login)
}
