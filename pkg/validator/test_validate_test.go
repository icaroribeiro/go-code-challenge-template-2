package validator_test

import (
	"fmt"
	"testing"

	validatorpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator"
	"github.com/stretchr/testify/assert"
)

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
