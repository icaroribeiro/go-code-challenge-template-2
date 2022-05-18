package route

import (
	"net/http"
)

// Route is the model of a route.
type Route struct {
	Name        string
	Method      string
	PathPrefix  string
	Path        string
	HandlerFunc http.HandlerFunc
}

// Routes is a slice of Route.
type Routes []Route
