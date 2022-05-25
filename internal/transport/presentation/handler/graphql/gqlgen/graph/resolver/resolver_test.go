package resolver_test

// import (
// 	"fmt"
// 	"testing"

// 	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/auth"
// 	resolverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/resolver
// 	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
// 	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/user"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// func TestResolverUnit(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }
// func (ts *TestSuite) TestNew() {

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			healthCheckService := new(healthcheckmockservice.Service)
// 			authService := new(authmockservice.Service)
// 			userService := new(usermockservice.Service)

// 			returnedResolver := resolverpkg.New(healthCheckService, authService, userService)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost")
// 			}
// 		})
// 	}
// }
