package auth_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/client"
	authservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/application/service/auth"
	userservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/application/service/user"
	healthcheckmockservice "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	persistententity "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/perentity"
	authdatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/repository/user"
	dbtrxdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/directive/dbtrx"
	authmockdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/mockdirective/auth"
	graphqlhandler "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/handler"
	authpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/auth"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/security"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestSignUp() {
	db := &gorm.DB{}

	authN := authpkg.New(ts.RSAKeys)

	credentials := securitypkg.CredentialsFactory(nil)

	opt := func(bd *client.Request) {}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningUp",
			SetUp: func(t *testing.T) {
				db = ts.DB

				opt = client.Var("input", credentials)
			},
			WantError: false,
			TearDown: func(t *testing.T) {
				result := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&persistententity.Auth{})
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
				result = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&persistententity.Login{})
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
				result = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&persistententity.User{})
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheDatabaseIsNull",
			SetUp: func(t *testing.T) {
				db = nil

				opt = client.Var("input", credentials)
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
			SetUp: func(t *testing.T) {
				db = ts.DB.Begin()

				opt = client.Var("input", credentials)
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(db)
			userDatastoreRepository := userdatastorerepository.New(db)
			loginDatastoreRepository := logindatastorerepository.New(db)

			healthCheckService := new(healthcheckmockservice.Service)
			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)

			userService := userservice.New(userDatastoreRepository, ts.Validator)

			dbTrxDirective := dbtrxdirective.New(db)

			authDirective := new(authmockdirective.Directive)
			authDirective.On("AuthMiddleware").Return(MockDirective())
			authDirective.On("AuthRenewalMiddleware").Return(MockDirective())

			graphqlHandler := graphqlhandler.New(healthCheckService, authService, userService, dbTrxDirective, authDirective)

			mutation := signUpMutation
			resp := SignUpMutationResponse{}

			srv := http.HandlerFunc(graphqlHandler.GraphQL())

			cl := client.New(srv)
			err := cl.Post(mutation, &resp, opt)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.NotEmpty(t, resp.SignUp.Token)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}

			tc.TearDown(t)
		})
	}
}
