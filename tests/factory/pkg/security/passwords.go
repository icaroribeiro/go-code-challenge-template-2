package security

import (
	"github.com/bluele/factory-go/factory"
	fake "github.com/brianvoe/gofakeit/v5"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/security"
	"golang.org/x/crypto/bcrypt"
)

// NewPasswords is the function that returns an instance of the user passwords for performing tests.
func NewPasswords(args map[string]interface{}) securitypkg.Passwords {
	passwordsFactory := factory.NewFactory(
		securitypkg.Passwords{},
	).Attr("CurrentPassword", func(fArgs factory.Args) (interface{}, error) {
		currentPassword := fake.Password(true, true, true, false, false, 8)

		security := securitypkg.New()
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

		security := securitypkg.New()
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

	return passwordsFactory.MustCreate().(securitypkg.Passwords)
}
