package customerror_test

import (
	"errors"
	"testing"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestNewf() {
	msg := ""
	err := errors.New("")

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingANewErrorWithFormattedMessage",
			SetUp: func(t *testing.T) {
				msg = "failed"
				err = customerror.Newf("%s", msg)
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedError := customerror.Newf("%s", msg)

			assert.Equal(t, err, returnedError)
		})
	}
}
