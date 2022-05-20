package tracer

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"
)

type (
	Tracer struct{}
)

var _ interface {
	graphql.ResponseInterceptor
} = Tracer{}

// ExtensionName is the function that...
func (a Tracer) ExtensionName() string {
	return "Tracer"
}

// Validate is the function that...
func (a Tracer) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptResponse is the function that...
func (a Tracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	log.Printf("%d before response\n", a)
	defer func() {
		log.Printf("%d after response\n", a)
	}()
	errList := graphql.GetErrors(ctx)
	if len(errList) > 0 {
		log.Println("Tem erros")
	} else {
		log.Println("Sem erros")
	}

	return next(ctx)
}
