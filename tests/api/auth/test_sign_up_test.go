package auth_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/client"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/auth"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/user"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	datastoreentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/entity"
	authdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/user"
	dbtrxdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/directive/dbtrx"
	authmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/mockdirective/auth"
	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/handler"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	securitypkgfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/pkg/security"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestSignUp() {
	db := &gorm.DB{}

	authN := authpkg.New(ts.RSAKeys)

	credentials := securitypkgfactory.NewCredentials(nil)

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
				result := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&datastoreentity.Auth{})
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
				result = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&datastoreentity.Login{})
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
				result = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&datastoreentity.User{})
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
