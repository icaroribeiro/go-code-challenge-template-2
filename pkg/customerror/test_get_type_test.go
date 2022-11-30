package customerror_test

import (
	"errors"
	"testing"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestGetType() {
	err := errors.New("")
	errorType := customerror.NoType

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingATypeWhenTheErrorHasAType",
			SetUp: func(t *testing.T) {
				err = customerror.BadRequest.New("failed")
				errorType = customerror.BadRequest
			},
		},
		{
			Context: "ItShouldSucceedInGettingNoTypeWhenTheErrorHasNoType",
			SetUp: func(t *testing.T) {
				err = errors.New("failed")
				errorType = customerror.NoType
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedErrorType := customerror.GetType(err)

			assert.Equal(t, errorType, returnedErrorType)
		})
	}
}
