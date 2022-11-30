package password

import (
	"reflect"
	"regexp"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	validatorv2 "gopkg.in/validator.v2"
)

// Validate is the function that validates a password.
func Validate(i interface{}, param string) error {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.String {
		return validatorv2.ErrUnsupported
	}

	if value.String() == "" {
		return nil
	}

	password := value.String()

	passwordRegex := regexp.MustCompile("[a-zA-Z0-9]{8,}")

	if !passwordRegex.MatchString(password) {
		return customerror.New("The password must contain only letters and digit with at least 8 characters")
	}

	return nil
}
