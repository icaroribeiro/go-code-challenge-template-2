package uuid

import (
	"reflect"
	"strings"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	_uuid "github.com/satori/go.uuid"
	validatorv2 "gopkg.in/validator.v2"
)

// Validate is the function that validates a UUID.
func Validate(i interface{}, param string) error {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.String {
		return validatorv2.ErrUnsupported
	}

	if strings.Compare(value.String(), "") == 0 {
		return nil
	}

	id := value.String()

	_, err := _uuid.FromString(id)
	if err != nil {
		return customerror.New("The id must contain an identifier following the UUID standard")
	}

	return nil
}
