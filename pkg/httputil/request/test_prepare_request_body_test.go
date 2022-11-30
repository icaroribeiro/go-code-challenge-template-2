package request_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestPrepareRequestBody() {
	var inputBody interface{}
	var reqBody io.Reader

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInPreparingRequestBodyIfInputBodyIsAnEmptyString",
			SetUp: func(t *testing.T) {
				inputBody = ""
				reqBody = nil
			},
		},
		{
			Context: "ItShouldSucceedInPreparingRequestBodyIfInputBodyIsAJsonStringWithEscapeSequencies",
			SetUp: func(t *testing.T) {
				inputBody = `
				{
					"testing":	"testing"
				}`
				reqBody = strings.NewReader(`{"testing":"testing"}`)
			},
		},
		{
			Context: "ItShouldSucceedInPreparingRequestBodyIfInputBodyIsAVariableSizedBufferOfBytes",
			SetUp: func(t *testing.T) {
				inputBody = new(bytes.Buffer)
				reqBody = new(bytes.Buffer)
			},
		},
		{
			Context: "ItShouldSucceedInPreparingRequestBodyIfInputBodyIsNeitherAStringNorAVariableSizedBufferOfBytes",
			SetUp: func(t *testing.T) {
				inputBody = nil
				reqBody = nil
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedBody := requesthttputilpkg.PrepareRequestBody(inputBody)

			assert.Equal(t, reqBody, returnedBody)
		})
	}
}
