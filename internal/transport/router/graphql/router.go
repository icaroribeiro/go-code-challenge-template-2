package swagger

import (
	"net/http"

	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
)

// ConfigureRoutes is the function that arranges the swagger's routes.
func ConfigureRoutes(swaggerHandler http.HandlerFunc, adapters map[string]adapterhttputilpkg.Adapter) routehttputilpkg.Routes {
	return routehttputilpkg.Routes{
		routehttputilpkg.Route{
			Name:       "Swagger",
			Method:     http.MethodGet,
			PathPrefix: "/swagger",
			HandlerFunc: adapterhttputilpkg.AdaptFunc(swaggerHandler).
				With(adapters["loggingMiddleware"]),
		},
	}
}
