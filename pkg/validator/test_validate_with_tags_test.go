package validator_test

import (
	"fmt"
	"testing"

	validatorpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator"
	uuidvalidatorpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator/uuid"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	validatorv2 "gopkg.in/validator.v2"
)

func (ts *TestSuite) TestValidateWithTags() {
	validationFuncs := map[string]validatorv2.ValidationFunc{
		"uuid": uuidvalidatorpkg.Validate,
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
