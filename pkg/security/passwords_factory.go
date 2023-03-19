package security

import (
	"github.com/bluele/factory-go/factory"
	fake "github.com/brianvoe/gofakeit/v5"
	"golang.org/x/crypto/bcrypt"
)

// PasswordsFactory is the function that returns an instance of the user passwords for performing tests.
func PasswordsFactory(args map[string]interface{}) Passwords {
	passwordsFactory := factory.NewFactory(
		Passwords{},
	).Attr("CurrentPassword", func(fArgs factory.Args) (interface{}, error) {
		currentPassword := fake.Password(true, true, true, false, false, 8)

		security := New()
		bytes, err := security.HashPassword(currentPassword, bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		hashedPassword := string(bytes)
		if val, ok := args["currentPassword"]; ok {
			hashedPassword = val.(string)
		}

		return hashedPassword, nil
	}).Attr("NewPassword", func(fArgs factory.Args) (interface{}, error) {
		newPassword := fake.Password(true, true, true, false, false, 8)

		security := New()
		bytes, err := security.HashPassword(newPassword, bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		hashedPassword := string(bytes)
		if val, ok := args["newPassword"]; ok {
			hashedPassword = val.(string)
		}

		return hashedPassword, nil
	})

	return passwordsFactory.MustCreate().(Passwords)
}
