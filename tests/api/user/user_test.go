package user_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/client"
	fake "github.com/brianvoe/gofakeit/v5"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/user"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/model"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/user"
	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql"
	authdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/auth"
	dbtrxmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/mockdirective/dbtrx"
	graphmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/model"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/adapter"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestUserInteg(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestGetAll() {
	dbTrx := &gorm.DB{}

	var authN authpkg.IAuth

	timeBeforeTokenExpTimeInSec := 30

	userDatastore := datastoremodel.User{}
	user := graphmodel.User{}
	loginDatastore := datastoremodel.Login{}
	authDatastore := datastoremodel.Auth{}

	key := ""
	bearerToken := []string{"", ""}
	value := ""

	adapters := map[string]adapterhttputilpkg.Adapter{
		"authMiddleware": authmiddlewarepkg.Auth(),
	}

	opt := func(bd *client.Request) {}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAllUsers",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				userDatastore = datastoremodel.User{
					Username: username,
				}

				result := dbTrx.Create(&userDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				domainUser := userDatastore.ToDomain()
				user.FromDomain(domainUser)

				loginDatastore = datastoremodel.Login{
					UserID:   userDatastore.ID,
					Username: username,
					Password: password,
				}

				result = dbTrx.Create(&loginDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				authDatastore = datastoremodel.Auth{
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

			healthCheckService := new(healthcheckmockservice.Service)
			authService := new(authmockservice.Service)

			userDatastoreRepository := userdatastorerepository.New(dbTrx)
			userService := userservice.New(userDatastoreRepository, ts.Validator)

			dbTrxDirective := new(dbtrxmockdirective.Directive)
			dbTrxDirective.On("DBTrxMiddleware").Return(MockDirective())

			authDirective := authdirective.New(dbTrx, authN, timeBeforeTokenExpTimeInSec)

			graphqlHandler := graphqlhandler.New(healthCheckService, authService, userService, dbTrxDirective, authDirective)

			query := getAllUsersQuery
			resp := GetAllUsersQueryResponse{}

			srv := AdaptHandlerWithHandlerFuncs(graphqlHandler.GraphQL(), adapters)

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
