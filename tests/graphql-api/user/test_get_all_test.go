package user_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/client"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/user"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	datastoreentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/entity"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/user"
	authdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/directive/auth"
	graphentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/entity"
	dbtrxmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/mockdirective/dbtrx"
	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/api/handler"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/adapter"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/auth"
	securitypkgfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/pkg/security"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestGetAll() {
	dbTrx := &gorm.DB{}

	authN := authpkg.New(ts.RSAKeys)

	credentials := securitypkgfactory.NewCredentials(nil)

	timeBeforeTokenExpTimeInSec := 30

	userDatastore := datastoreentity.User{}
	user := graphentity.User{}
	loginDatastore := datastoreentity.Login{}
	authDatastore := datastoreentity.Auth{}

	key := ""
	bearerToken := []string{"", ""}
	value := ""

	adapters := map[string]adapterhttputilpkg.Adapter{
		"authMiddleware": authmiddlewarepkg.Auth(authN),
	}

	opt := func(bd *client.Request) {}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAllUsers",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				userDatastore = datastoreentity.User{
					Username: credentials.Username,
				}

				result := dbTrx.Create(&userDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				domainUser := userDatastore.ToDomain()
				user.FromDomain(domainUser)

				loginDatastore = datastoreentity.Login{
					UserID:   userDatastore.ID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				result = dbTrx.Create(&loginDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				authDatastore = datastoreentity.Auth{
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
			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			userDatastoreRepository := userdatastorerepository.New(dbTrx)

			healthCheckService := new(healthcheckmockservice.Service)
			authService := new(authmockservice.Service)
			userService := userservice.New(userDatastoreRepository, ts.Validator)

			dbTrxDirective := new(dbtrxmockdirective.Directive)
			dbTrxDirective.On("DBTrxMiddleware").Return(MockDirective())

			authDirective := authdirective.New(dbTrx, authN, timeBeforeTokenExpTimeInSec)

			graphqlHandler := graphqlhandler.New(healthCheckService, authService, userService, dbTrxDirective, authDirective)

			query := getAllUsersQuery
			resp := GetAllUsersQueryResponse{}

			srv := http.HandlerFunc(adapterhttputilpkg.AdaptFunc(graphqlHandler.GraphQL()).
				With(adapters["authMiddleware"]))

			cl := client.New(srv)
			err := cl.Post(query, &resp, opt)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, user.ID.String(), resp.GetAllUsers[0].ID)
				assert.Equal(t, user.Username, resp.GetAllUsers[0].Username)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}

			tc.TearDown(t)
		})
	}
}
