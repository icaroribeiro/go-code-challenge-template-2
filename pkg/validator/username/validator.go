package username

import (
	"reflect"
	"regexp"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	validatorv2 "gopkg.in/validator.v2"
)

// Validate is the function that validates a username.
func Validate(i interface{}, param string) error {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.String {
		return validatorv2.ErrUnsupported
	}

	if value.String() == "" {
		return nil
	}

	username := value.String()

	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9]{5,}$")

	if !usernameRegex.MatchString(username) {
		return customerror.New("The username must contain only letters and digit with at least 5 characters")
	}

	return nil
}
