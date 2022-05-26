package resolver_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/DATA-DOG/go-sqlmock"
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/dgrijalva/jwt-go"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
	authdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/generated"
	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/resolver"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainfactorymodel "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/model"
	datastorefactorymodel "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/infrastructure/storage/datastore/model"
	mockauthpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/mocks/pkg/mockauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestUserResolversUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestGetAllUsers() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	ctx := context.Background()

	tokenString := fake.Word()

	token := &jwt.Token{}

	auth := domainfactorymodel.NewAuth(nil)

	sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

	args := map[string]interface{}{
		"id":     auth.ID,
		"userID": auth.UserID,
	}

	authDatastore := datastorefactorymodel.NewAuth(args)

	user := domainmodel.User{}

	dbTrx := &gorm.DB{}
	dbTrx = nil

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAllUsers",
			SetUp: func(t *testing.T) {
				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(authDatastore.ID).
					WillReturnRows(rows)

				user = domainfactorymodel.NewUser(nil)

				returnArgs = ReturnArgs{
					{token, nil},
					{auth, nil},
					{domainmodel.Users{user}, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenGettingAllUsers",
			SetUp: func(t *testing.T) {
				rows := sqlmock.
					NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(authDatastore.ID, authDatastore.UserID, authDatastore.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(authDatastore.ID).
					WillReturnRows(rows)

				returnArgs = ReturnArgs{
					{token, nil},
					{auth, nil},
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
			userService.On("GetAll").Return(returnArgs[2]...)

			resolver := resolverpkg.New(healthCheckService, authService, userService)

			c := generated.Config{Resolvers: resolver}

			authN := new(mockauthpkg.Auth)
			authN.On("DecodeToken", tokenString).Return(returnArgs[0]...)
			authN.On("FetchAuthFromToken", token).Return(returnArgs[1]...)

			c.Directives.UseAuthMiddleware = authdirectivepkg.AuthMiddleware(db, authN)

			srv := handler.NewDefaultServer(
				generated.NewExecutableSchema(
					c,
				),
			)

			cl := client.New(srv)

			query := getAllUsersQuery

			resp := GetAllUsersResponse{}

			err := cl.Post(query, &resp, addTokenStringCtxValue(ctx, tokenString))

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, user.ID.String(), resp.GetAllUsers[0].ID)
				assert.Equal(t, user.Username, resp.GetAllUsers[0].Username)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}
