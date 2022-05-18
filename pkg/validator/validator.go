package validator

import (
	"fmt"

	validatorv2 "gopkg.in/validator.v2"
)

type Validator struct {
	ValidatorV2 *validatorv2.Validator
}

// New is the factory function that encapsulates the implementation related to validator.
func New(validationFuncs map[string]validatorv2.ValidationFunc) (IValidator, error) {
	vv2 := validatorv2.NewValidator()

	for name, validationFunc := range validationFuncs {
		err := vv2.SetValidationFunc(name, validationFunc)
		if err != nil {
			msg := "failed to set the validation function for the %s field: %s"
			return &Validator{}, fmt.Errorf(msg, name, err.Error())
		}
	}

	return &Validator{
		ValidatorV2: vv2,
	}, nil
}

// Validateis the function that validates the fields of structs based on 'validator' tags.
func (v *Validator) Validate(i interface{}) error {
	return v.ValidatorV2.Validate(i)
}

// ValidateWithTags is the function that validates a value based on the provided tags.
func (v *Validator) ValidateWithTags(i interface{}, tags string) error {
	return v.ValidatorV2.Valid(i, tags)
}
