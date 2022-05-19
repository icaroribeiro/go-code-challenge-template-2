package adapter_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"

	"github.com/gorilla/mux"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/adapter"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/response"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestAdapterUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestAdaptFunc() {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInAdaptingAFunction",
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			returnedAdaptedHandlerFunc := adapterhttputilpkg.AdaptFunc(handlerFunc)

			handlerFunc1 := runtime.FuncForPC(reflect.ValueOf(handlerFunc).Pointer()).Name()
			handlerFunc2 := runtime.FuncForPC(reflect.ValueOf(returnedAdaptedHandlerFunc.HandlerFunc).Pointer()).Name()
			assert.Equal(t, handlerFunc1, handlerFunc2)
		})
	}
}

var countCtxKey = &contextKey{"count"}

type contextKey struct {
	name string
}

func addOneToCountHandlerFunc() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			count := 0
			ok := false

			ctx := r.Context()
			i := ctx.Value(countCtxKey)
			if i == nil {
				count = 1
			} else {
				count, ok = i.(int)
				if !ok {
					responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
					return
				}
				count += 1
			}

			ctx = context.WithValue(ctx, countCtxKey, count)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	}
}

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
