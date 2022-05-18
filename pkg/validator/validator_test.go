package validator_test

import (
	"fmt"
	"testing"

	validatorpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/validator"
	uuidvalidator "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/validator/uuid"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	validatorv2 "gopkg.in/validator.v2"
)

func TestValidatorUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestNew() {
	validationFuncs := make(map[string]validatorv2.ValidationFunc)

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInInitializingTheValidator",
			SetUp: func(t *testing.T) {
				validationFuncs = make(map[string]validatorv2.ValidationFunc)
				validationFuncs["foo"] = func(i interface{}, s string) error {
					return nil
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheValidationFunctionDoesNotHaveAName",
			SetUp: func(t *testing.T) {
				validationFuncs = make(map[string]validatorv2.ValidationFunc)
				validationFuncs[""] = func(i interface{}, s string) error {
					return nil
				}
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			_, err := validatorpkg.New(validationFuncs)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost")
			}
		})
	}
}

func (ts *TestSuite) TestValidate() {
	validator, err := validatorpkg.New(nil)
	assert.Nil(ts.T(), err, fmt.Sprintf("Unexpected error: %v", err))

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedWithTheValidateTag",
			Inf: Foo{
				Field1: "field1",
			},
			WantError: false,
		},
		{
			Context:   "ItShouldFailWithTheValidateTagBecuaseField1IsEmpty",
			Inf:       Foo{},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			err := validator.Validate(tc.Inf)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost")
			}
		})
	}
}

func (ts *TestSuite) TestValidateWithTags() {
	validationFuncs := map[string]validatorv2.ValidationFunc{
		"uuid": uuidvalidator.Validate,
	}

	validator, err := validatorpkg.New(validationFuncs)
	assert.Nil(ts.T(), err, fmt.Sprintf("Unexpected error: %v", err))

	ts.Cases = Cases{
		{
			Context:   "ItShouldSucceed",
			Inf:       uuid.NewV4().String(),
			Tags:      "uuid",
			WantError: false,
		},
		{
			Context:   "ItShouldFailIfTheTagIsEmpty",
			Inf:       "",
			Tags:      "",
			WantError: true,
		},
		{
			Context:   "ItShouldFailIfTheTagDoesNotExist",
			Inf:       "",
			Tags:      "foo",
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			err := validator.ValidateWithTags(tc.Inf, tc.Tags)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost")
			}
		})
	}
}
