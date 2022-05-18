package handler

import "net/http"

// GetNotFoundHandler is the function that sets handler to be used when no route matches.
func GetNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}

// GetMethodNotAllowedHandler is the function that sets handler to be used when the request method does not match the route.
func GetMethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}
