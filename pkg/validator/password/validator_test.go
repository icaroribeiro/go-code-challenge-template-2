package password_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	passwordvalidatorpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/validator/password"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestPasswordUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestValidate() {
	ts.Cases = Cases{
		{
			Context:   "ItShouldSucceedInNotValidatingPasswordIfTheStringIsEmpty",
			Inf:       "",
			Param:     "",
			WantError: false,
		},
		{
			Context:   "ItShouldSucceedIfTheStringContainsAValidPassword",
			Inf:       fake.Password(true, true, true, false, false, 8),
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
			Context:   "ItShouldFailIfTheStringContainsLessThanEightCharacters",
			Inf:       fake.Password(true, true, true, false, false, 7),
			Param:     "",
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			err := passwordvalidatorpkg.Validate(tc.Inf, tc.Param)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost")
			}
		})
	}
}
