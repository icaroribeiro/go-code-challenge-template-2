package graphql_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gorilla/mux"
// 	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/mockservice/healthcheck"
// 	healthcheckhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/handler/healthcheck"
// 	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
// 	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
// 	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
// 	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// func TestHandlerUnit(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestGetStatus() {
// 	text := "everything is up and running"

// 	message := messagehttputilpkg.Message{Text: text}

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInGettingStatus",
// 			SetUp: func(t *testing.T) {
// 				returnArgs = ReturnArgs{
// 					{nil},
// 				}
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 		},
// 		{
// 			Context: "ItShouldFailIfItAnErrorOccursWhenGettingTheStatus",
// 			SetUp: func(t *testing.T) {
// 				returnArgs = ReturnArgs{
// 					{customerror.New("failed")},
// 				}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			healthCheckService := new(healthcheckmockservice.Service)
// 			healthCheckService.On("GetStatus").Return(returnArgs[0]...)

// 			healthCheckHandler := healthcheckhandler.New(healthCheckService)

// 			route := routehttputilpkg.Route{
// 				Name:        "GetStatus",
// 				Method:      http.MethodGet,
// 				Path:        "/status",
// 				HandlerFunc: healthCheckHandler.GetStatus,
// 			}

// 			requestData := requesthttputilpkg.RequestData{
// 				Method: route.Method,
// 				Target: route.Path,
// 			}

// 			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

// 			resprec := httptest.NewRecorder()

// 			router := mux.NewRouter()

// 			router.Name(route.Name).
// 				Methods(route.Method).
// 				Path(route.Path).
// 				HandlerFunc(route.HandlerFunc)

// 			router.ServeHTTP(resprec, req)

// 			if !tc.WantError {
// 				assert.Equal(t, resprec.Code, tc.StatusCode)
// 				returnedMessage := messagehttputilpkg.Message{}
// 				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
// 				assert.Equal(t, returnedMessage.Text, message.Text)
// 			} else {
// 				assert.Equal(t, resprec.Code, tc.StatusCode)
// 			}
// 		})
// 	}
// }
