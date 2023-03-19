package auth_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/client"
	authservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/application/service/auth"
	userservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/application/service/user"
	healthcheckmockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	persistententity "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/perentity"
	authdatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/repository/user"
	authdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/directive/auth"
	dbtrxmockdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/mockdirective/dbtrx"
	graphqlhandler "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/handler"
	authpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/auth"
	adapterhttputilpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/httputil/adapter"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/middleware/auth"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/security"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestRefreshToken() {
	dbTrx := &gorm.DB{}

	authN := authpkg.New(ts.RSAKeys)

	credentials := securitypkg.CredentialsFactory(nil)

	timeBeforeTokenExpTimeInSec := 120

	userDatastore := persistententity.User{}
	loginDatastore := persistententity.Login{}
	authDatastore := persistententity.Auth{}

	key := ""
	bearerToken := []string{"", ""}
	value := ""

	adapters := map[string]adapterhttputilpkg.Adapter{
		"authMiddleware": authmiddlewarepkg.Auth(authN),
	}

	opt := func(bd *client.Request) {}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRefreshingTheToken",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				userDatastore = persistententity.User{
					Username: credentials.Username,
				}

				result := dbTrx.Create(&userDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				loginDatastore = persistententity.Login{
					UserID:   userDatastore.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				result = dbTrx.Create(&loginDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				authDatastore = persistententity.Auth{
					UserID: userDatastore.ID,
				}

				result = dbTrx.Create(&authDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				key = "Authorization"
				tokenString, err := authN.CreateToken(authDatastore.ToDomain(), ts.TokenExpTimeInSec)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				bearerToken = []string{"Bearer", tokenString}
				value = strings.Join(bearerToken[:], " ")

				opt = AddRequestHeaderEntries(key, value)
			},
			WantError: false,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheTokenIsNotSentIntheRequest",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))
			},
			WantError: true,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(dbTrx)
			loginDatastoreRepository := logindatastorerepository.New(dbTrx)
			userDatastoreRepository := userdatastorerepository.New(dbTrx)

			healthCheckService := new(healthcheckmockservice.Service)
			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
			userService := userservice.New(userDatastoreRepository, ts.Validator)

			dbTrxDirective := new(dbtrxmockdirective.Directive)
			dbTrxDirective.On("DBTrxMiddleware").Return(MockDirective())

			authDirective := authdirective.New(dbTrx, authN, timeBeforeTokenExpTimeInSec)

			graphqlHandler := graphqlhandler.New(healthCheckService, authService, userService, dbTrxDirective, authDirective)

			mutation := refreshTokenMutation
			resp := RefreshTokenMutationResponse{}

			srv := http.HandlerFunc(adapterhttputilpkg.AdaptFunc(graphqlHandler.GraphQL()).
				With(adapters["authMiddleware"]))

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, opt)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.NotEmpty(t, resp.RefreshToken.Token)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}

			tc.TearDown(t)
		})
	}
}
