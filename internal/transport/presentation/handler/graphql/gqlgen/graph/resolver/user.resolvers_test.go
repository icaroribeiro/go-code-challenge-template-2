package resolver_test

import (
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/resolver"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestUserResolversUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestGetAllUsers() {
	user := domainmodel.User{}

	dbTrx := &gorm.DB{}
	dbTrx = nil

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAllUsers",
			SetUp: func(t *testing.T) {
				user = domainmodelfactory.NewUser(nil)

				returnArgs = ReturnArgs{
					{domainmodel.Users{user}, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenGettingAllUsers",
			SetUp: func(t *testing.T) {
				returnArgs = ReturnArgs{
					{domainmodel.Users{}, customerror.New("failed")},
				}
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			healthCheckService := new(healthcheckmockservice.Service)
			authService := new(authmockservice.Service)

			userService := new(usermockservice.Service)
			userService.On("WithDBTrx", dbTrx).Return(userService)
			userService.On("GetAll").Return(returnArgs[0]...)

			resolver := resolverpkg.New(healthCheckService, authService, userService)

			c := generated.Config{Resolvers: resolver}

			c.Directives.UseAuthMiddleware = MockSchemaDirective()

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					c,
				),
			)

			query := getAllUsersQuery
			resp := GetAllUsersQueryResponse{}

			cl := client.New(srv)
			err := cl.Post(query, &resp)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, user.ID.String(), resp.GetAllUsers[0].ID)
				assert.Equal(t, user.Username, resp.GetAllUsers[0].Username)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}
		})
	}
}
