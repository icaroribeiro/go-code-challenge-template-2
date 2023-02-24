package uuid

import (
	"reflect"

	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	_uuid "github.com/satori/go.uuid"
	validatorv2 "gopkg.in/validator.v2"
)

// Validate is the function that validates a UUID.
func Validate(i interface{}, param string) error {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.String {
		return validatorv2.ErrUnsupported
	}

	if id, _ := _uuid.FromString(value.String()); id == _uuid.Nil {
		return customerror.New("The id must contain an identifier following the UUID standard not empty")
	}

	return nil
}
