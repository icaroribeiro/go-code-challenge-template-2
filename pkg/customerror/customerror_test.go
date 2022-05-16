package customerror_test

import (
	"errors"
	"testing"

	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestCustomErrorUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

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

func (ts *TestSuite) TestNewfTypedError() {
	errorType := customerror.NoType
	msg := ""
	err := errors.New("")

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingANewTypedErrorWithFormattedMessage",
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

			returnedError := errorType.Newf("%s", msg)

			assert.Equal(t, err, returnedError)
		})
	}
}

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
