package customerror_test

import (
	"errors"
	"testing"

	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestNewfTypedError() {
	errorType := customerror.NoType
	msg := ""
	err := errors.New("")

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingANewfTypedErrorWithMessage",
			SetUp: func(t *testing.T) {
				errorType = customerror.BadRequest
				msg = "failed"
				err = customerror.BadRequest.Newf("%s", msg)
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedError := errorType.Newf(msg)

			assert.Equal(t, err, returnedError)
		})
	}
}
