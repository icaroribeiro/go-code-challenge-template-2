package customerror_test

import (
	"errors"
	"testing"

	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestNew() {
	msg := ""
	err := errors.New("")

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingANewErrorWithMessage",
			SetUp: func(t *testing.T) {
				msg = "failed"
				err = customerror.New(msg)
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedError := customerror.New(msg)

			assert.Equal(t, err, returnedError)
		})
	}
}
