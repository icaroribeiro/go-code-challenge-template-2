package security

import (
	"github.com/bluele/factory-go/factory"
	fake "github.com/brianvoe/gofakeit/v5"
	"golang.org/x/crypto/bcrypt"
)

// CredentialsFactory is the function that returns an instance of the user credentials for performing tests.
func CredentialsFactory(args map[string]interface{}) Credentials {
	credentialsFactory := factory.NewFactory(
		Credentials{},
	).Attr("Username", func(fArgs factory.Args) (interface{}, error) {
		username := fake.Username()
		if val, ok := args["username"]; ok {
			username = val.(string)
		}

		return username, nil
	}).Attr("Password", func(fArgs factory.Args) (interface{}, error) {
		password := fake.Password(true, true, true, false, false, 8)

		security := New()
		bytes, err := security.HashPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		hashedPassword := string(bytes)
		if val, ok := args["password"]; ok {
			hashedPassword = val.(string)
		}

		return hashedPassword, nil
	})

	return credentialsFactory.MustCreate().(Credentials)
}
