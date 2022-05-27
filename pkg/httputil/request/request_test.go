package request_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestRequestUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

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

func (ts *TestSuite) TestSetRequestHeaders() {
	key := ""
	value := ""
	headers := make(map[string][]string)

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInNotSettingRequestHeadersIfHeadersMapIsEmpty",
			SetUp: func(t *testing.T) {
				key = "Content-Type"
				headers = map[string][]string{}
			},
		},
		{
			Context: "ItShouldSucceedInSettingRequestHeaders",
			SetUp: func(t *testing.T) {
				key = "Content-Type"
				value = "application/json"
				headers = map[string][]string{
					key: {value},
				}
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			route := routehttputilpkg.Route{
				Name:   "Testing",
				Method: http.MethodGet,
				Path:   "/testing",
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			requesthttputilpkg.SetRequestHeaders(req, headers)

			sort.Strings(headers[key])
			sort.Strings(req.Header.Values(key))
			assert.Equal(t, headers[key], req.Header.Values(key))
		})
	}
}
