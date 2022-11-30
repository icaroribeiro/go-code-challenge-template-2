package customerror_test

import (
	"errors"
	"testing"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestNewTypedError() {
	errorType := customerror.NoType
	msg := ""
	err := errors.New("")

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingANewTypedErrorWithMessage",
			SetUp: func(t *testing.T) {
				errorType = customerror.BadRequest
				msg = "failed"
				err = customerror.BadRequest.New(msg)
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedError := errorType.New(msg)

			assert.Equal(t, err, returnedError)
		})
	}
}
