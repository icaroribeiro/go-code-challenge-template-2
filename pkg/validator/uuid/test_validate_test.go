package uuid_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	uuidvalidatorpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator/uuid"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestValidate() {
	ts.Cases = Cases{
		{
			Context:   "ItShouldSucceedIfTheStringIsAValidUUID",
			Inf:       uuid.NewV4().String(),
			Param:     "",
			WantError: false,
		},
		{
			Context:   "ItShouldSucceedInNotValidatingUUIDIfTheStringIsEmpty",
			Inf:       "",
			Param:     "",
			WantError: false,
		},
		{
			Context:   "ItShouldFailIfANewValueIsNotInitializedToTheConcreteValue",
			Inf:       nil,
			Param:     "",
			WantError: true,
		},
		{
			Context:   "ItShouldSucceedIfTheStringIsNotAUUID",
			Inf:       fake.Word(),
			Param:     "",
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			err := uuidvalidatorpkg.Validate(tc.Inf, tc.Param)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost")
			}
		})
	}
}
