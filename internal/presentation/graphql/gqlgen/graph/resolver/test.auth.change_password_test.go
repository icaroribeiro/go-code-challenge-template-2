package resolver_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/99designs/gqlgen/client"
// 	"github.com/99designs/gqlgen/graphql/handler"
// 	domainentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/entity"
// 	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
// 	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
// 	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
// 	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/generated"
// 	authmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/mockdirective/auth"
// 	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/resolver"
// 	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
// 	domainentityfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/entity"
// 	securitypkgfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/pkg/security"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/gorm"
// )

// func (ts *TestSuite) TestChangePassword() {
// 	dbTrx := &gorm.DB{}
// 	dbTrx = nil

// 	passwords := securitypkgfactory.NewPasswords(nil)

// 	auth := domainentityfactory.NewAuth(nil)

// 	opts := []client.Option{}

// 	message := ""

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInResettingThePassword",
// 			SetUp: func(t *testing.T) {
// 				opts = []client.Option{}
// 				opts = append(opts, client.Var("input", passwords))
// 				ctx := context.Background()
// 				opts = append(opts, AddAuthDetailsToCtx(ctx, auth))

// 				message = "the password has been updated successfully"

// 				returnArgs = ReturnArgs{
// 					{nil},
// 				}
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
// 			SetUp: func(t *testing.T) {
// 				opts = []client.Option{}
// 				opts = append(opts, client.Var("input", passwords))

// 				returnArgs = ReturnArgs{
// 					{nil},
// 				}
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfItTheAuthFromTheRequestContextIsEmpty",
// 			SetUp: func(t *testing.T) {
// 				opts = []client.Option{}
// 				opts = append(opts, client.Var("input", passwords))
// 				ctx := context.Background()
// 				opts = append(opts, AddAuthDetailsToCtx(ctx, domainentity.Auth{}))

// 				returnArgs = ReturnArgs{
// 					{nil},
// 				}
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenResettingThePassword",
// 			SetUp: func(t *testing.T) {
// 				opts = []client.Option{}
// 				opts = append(opts, client.Var("input", passwords))
// 				ctx := context.Background()
// 				opts = append(opts, AddAuthDetailsToCtx(ctx, auth))

// 				returnArgs = ReturnArgs{
// 					{customerror.New("failed")},
// 				}
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			healthCheckService := new(healthcheckmockservice.Service)
// 			authService := new(authmockservice.Service)
// 			authService.On("WithDBTrx", dbTrx).Return(authService)
// 			authService.On("ModifyPassword", auth.UserID.String(), passwords).Return(returnArgs[0]...)
// 			userService := new(usermockservice.Service)

// 			res := resolverpkg.New(healthCheckService, authService, userService)

// 			cfg := generated.Config{Resolvers: res}

// 			authDirective := new(authmockdirective.Directive)
// 			authDirective.On("AuthMiddleware").Return(MockDirective())

// 			cfg.Directives.UseAuthMiddleware = authDirective.AuthMiddleware()

// 			srv := handler.NewDefaultServer(
// 				generated.NewExecutableSchema(
// 					cfg,
// 				),
// 			)

// 			mutation := changePasswordMutation
// 			resp := ChangePasswordMutationResponse{}

// 			cl := client.New(srv)
// 			err := cl.Post(mutation, &resp, opts...)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
// 				assert.NotEmpty(t, resp.ChangePassword.Message)
// 				assert.Equal(t, message, resp.ChangePassword.Message)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 			}
// 		})
// 	}
// }
