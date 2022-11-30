package adapter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestWith() {
	expectedCount := 0
	statusCode := 0
	handlerFuncs := []adapterhttputilpkg.Adapter{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInAdaptingAddOneToCountHandlerFuncOnce",
			SetUp: func(t *testing.T) {
				expectedCount = 1
				statusCode = http.StatusOK
				handlerFuncs = []adapterhttputilpkg.Adapter{addOneToCountHandlerFunc()}
			},
		},
		{
			Context: "ItShouldSucceedInAdaptingAddOneToCountHandlerFuncTwice",
			SetUp: func(t *testing.T) {
				expectedCount = 2
				statusCode = http.StatusOK
				handlerFuncs = []adapterhttputilpkg.Adapter{addOneToCountHandlerFunc(), addOneToCountHandlerFunc()}
			},
		},
		{
			Context: "ItShouldSucceedInAdaptingAddOneToCountHandlerFuncThrice",
			SetUp: func(t *testing.T) {
				expectedCount = 3
				statusCode = http.StatusOK
				handlerFuncs = []adapterhttputilpkg.Adapter{
					addOneToCountHandlerFunc(), addOneToCountHandlerFunc(), addOneToCountHandlerFunc(),
				}
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			countHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
				i := r.Context().Value(countCtxKey)

				count, ok := i.(int)
				if !ok {
					responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
					return
				}

				assert.Equal(t, expectedCount, count)
			}

			returnedHandlerFunc := adapterhttputilpkg.AdaptFunc(countHandlerFunc).With(handlerFuncs...)

			route := routehttputilpkg.Route{
				Name:        "Testing",
				Method:      http.MethodPost,
				Path:        "/testing",
				HandlerFunc: returnedHandlerFunc,
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			assert.Equal(t, statusCode, resprec.Result().StatusCode)
		})
	}
}
