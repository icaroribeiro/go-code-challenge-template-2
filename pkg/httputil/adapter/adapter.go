package adapter

import (
	"net/http"
)

// AdaptedHandlerFunc is a middleware handler that holds an http.HandlerFunc.
type AdaptedHandlerFunc struct {
	HandlerFunc http.HandlerFunc
}

// Adapter is the function that both takes in and returns an http.HandlerFunc.
type Adapter func(http.HandlerFunc) http.HandlerFunc

// AdaptFunc is the function that receives an http.HandlerFunc in order to be adapted with other functions.
func AdaptFunc(f func(w http.ResponseWriter, r *http.Request)) AdaptedHandlerFunc {
	return AdaptedHandlerFunc{
		HandlerFunc: f,
	}
}

// With is the function that applies a chain of http.HandlerFuncs (adapters) to an http.HandlerFunc.
func (a AdaptedHandlerFunc) With(adapters ...Adapter) http.HandlerFunc {
	handlerFunc := a.HandlerFunc

	lastPos := len(adapters) - 1

	for index := range adapters {
		handlerFunc = adapters[lastPos-index](handlerFunc)
	}

	return handlerFunc
}
