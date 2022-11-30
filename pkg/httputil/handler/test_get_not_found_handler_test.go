package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	handlerhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/handler"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestGetNotFoundHandler() {
	statusCode := 0

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInConfiguringNotFoundHandler",
			SetUp: func(t *testing.T) {
				statusCode = http.StatusNotFound
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

			resprec := httptest.NewRecorder()

			handler := handlerhttputilpkg.GetNotFoundHandler()

			handler.ServeHTTP(resprec, req)

			assert.Equal(t, statusCode, resprec.Result().StatusCode)
		})
	}
}
