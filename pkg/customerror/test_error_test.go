package customerror_test

import (
	"errors"
	"testing"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestError() {
	msg := "failed"
	err := errors.New(msg)

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInReturningTheErrorMessage",
			SetUp: func(t *testing.T) {
				err = customerror.New(msg)
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedMsg := err.Error()

			assert.Equal(t, msg, returnedMsg)
		})
	}
}
