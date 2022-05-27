package user_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	fake "github.com/brianvoe/gofakeit/v5"
// 	"github.com/gorilla/mux"
// 	userservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/application/service/user"
// 	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/model"
// 	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/user"
// 	userhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/handler/user"
// 	presentationmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/model"
// 	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
// 	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/gorm"
// )

// func TestUserInteg(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestGetAll() {
// 	dbTrx := &gorm.DB{}

// 	userDatastore := datastoremodel.User{}

// 	user := presentationmodel.User{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInGettingAllUsers",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				username := fake.Username()

// 				userDatastore = datastoremodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&userDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				domainUser := userDatastore.ToDomain()
// 				user.FromDomain(domainUser)
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			userDatastoreRepository := userdatastorerepository.New(dbTrx)
// 			userService := userservice.New(userDatastoreRepository, ts.Validator)
// 			userHandler := userhandler.New(userService)

// 			route := routehttputilpkg.Route{
// 				Name:        "GetAllUsers",
// 				Method:      "GET",
// 				Path:        "/users",
// 				HandlerFunc: userHandler.GetAll,
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
// 				returnedUsers := presentationmodel.Users{}
// 				err := json.NewDecoder(resprec.Body).Decode(&returnedUsers)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
// 				assert.Equal(t, user.ID, returnedUsers[0].ID)
// 				assert.Equal(t, user.Username, returnedUsers[0].Username)
// 			} else {
// 				assert.Equal(t, resprec.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }
