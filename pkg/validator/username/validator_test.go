package username_test

import (
	"fmt"
	"testing"

	usernamevalidator "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/validator/username"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestUsernameUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestValidate() {
	ts.Cases = Cases{
		{
			Context:   "ItShouldSucceedInNotValidatingUsernameIfTheStringIsEmpty",
			Inf:       "",
			Param:     "",
			WantError: false,
		},
		{
			Context:   "ItShouldSucceedIfTheStringContainsAValidUsername",
			Inf:       "foobar123",
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
			Context:   "ItShouldFailIfTheStringDoesNotContainOnlyLettersAndNumbers",
			Inf:       "foo.bar123",
			Param:     "",
			WantError: true,
		},
		{
			Context:   "ItShouldFailIfTheStringContainsLessThanFiveCharacters",
			Inf:       "foo",
			Param:     "",
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			err := usernamevalidator.Validate(tc.Inf, tc.Param)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost")
			}
		})
	}
}
