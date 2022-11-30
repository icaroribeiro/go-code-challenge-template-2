package validator_test

import (
	"fmt"
	"testing"

	validatorpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator"
	"github.com/stretchr/testify/assert"
	validatorv2 "gopkg.in/validator.v2"
)

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
