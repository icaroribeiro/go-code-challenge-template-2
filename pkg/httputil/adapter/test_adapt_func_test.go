package adapter_test

import (
	"context"
	"net/http"
	"reflect"
	"runtime"
	"testing"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
	"github.com/stretchr/testify/assert"
)

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
