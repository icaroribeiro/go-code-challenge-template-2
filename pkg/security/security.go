package security

import (
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"golang.org/x/crypto/bcrypt"
)

type Security struct{}

// New is the factory function that encapsulates the implementation related to security.
func New() ISecurity {
	return &Security{}
}

// HashPassword is the function that encrypts a password using the bcrypt algorithm and a default cost of hashing.
func (s *Security) HashPassword(password string, cost int) ([]byte, error) {
	if password == "" {
		return []byte{}, customerror.BadRequest.New("the password is empty")
	}

	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

// VerifyPasswords is the function that verifies if a hashed password is equivalent to a plaintext.
func (s *Security) VerifyPasswords(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return customerror.Unauthorized.New("the password is invalid")
		}

		return err
	}

	return nil
}
