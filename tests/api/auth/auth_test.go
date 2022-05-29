package auth_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/client"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/auth"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/user"
	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/application/mockservice/healthcheck"
	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/model"
	authdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/user"
	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql"
	dbtrxdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/dbtrx"
	authmockdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/mockdirective/auth"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	securitypkgfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestAuthInteg(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestSignUp() {
	db := &gorm.DB{}

	var authN authpkg.IAuth

	credentials := securitypkgfactory.NewCredentials(nil)

	opt := func(bd *client.Request) {}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningUp",
			SetUp: func(t *testing.T) {
				db = ts.DB

				authN = authpkg.New(ts.RSAKeys)

				opt = client.Var("input", credentials)
			},
			WantError: false,
			TearDown: func(t *testing.T) {
				result := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&datastoremodel.Auth{})
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
				result = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&datastoremodel.Login{})
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
				result = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&datastoremodel.User{})
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheDatabaseIsNull",
			SetUp: func(t *testing.T) {
				db = nil

				authN = authpkg.New(ts.RSAKeys)

				opt = client.Var("input", credentials)
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
			SetUp: func(t *testing.T) {
				db = ts.DB.Begin()

				authN = authpkg.New(ts.RSAKeys)

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

// func (ts *TestSuite) TestSignIn() {
// 	dbTrx := &gorm.DB{}

// 	var authN authpkg.IAuth

// 	userDatastore := datastoremodel.User{}

// 	loginDatastore := datastoremodel.Login{}

// 	credentials := securitypkg.Credentials{}

// 	body := ""

// 	dbTrxCtxValue := &gorm.DB{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSigningIn",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				userDatastore = datastoremodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&userDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				loginDatastore = datastoremodel.Login{
// 					UserID:   userDatastore.ID,
// 					Username: username,
// 					Password: password,
// 				}

// 				result = dbTrx.Create(&loginDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				credentials = securitypkg.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"username":"%s",
// 					"password":"%s"
// 				}`,
// 					credentials.Username, credentials.Password)

// 				dbTrxCtxValue = dbTrx
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
// 			SetUp: func(t *testing.T) {
// 				authN = authpkg.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = securitypkg.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = ""

// 				dbTrxCtxValue = nil
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = securitypkg.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = fmt.Sprintf(`
// 					"username":"%s",
// 					"password":"%s"
// 				`,
// 					credentials.Username, credentials.Password)

// 				dbTrxCtxValue = dbTrx
// 			},
// 			StatusCode: http.StatusBadRequest,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = securitypkg.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"username":"%s",
// 					"password":"%s"
// 				}`,
// 					credentials.Username, credentials.Password)

// 				dbTrxCtxValue = dbTrx
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDatastoreRepository := authdatastorerepository.New(dbTrx)

// 			userDatastoreRepository := userdatastorerepository.New(dbTrx)

// 			loginDatastoreRepository := logindatastorerepository.New(dbTrx)

// 			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
// 				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
// 			authHandler := authhandler.New(authService)

// 			route := routehttputilpkg.Route{
// 				Name:        "SignIn",
// 				Method:      http.MethodPost,
// 				Path:        "/sign_in",
// 				HandlerFunc: authHandler.SignIn,
// 			}

// 			requestData := requesthttputilpkg.RequestData{
// 				Method: route.Method,
// 				Target: route.Path,
// 				Body:   body,
// 			}

// 			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

// 			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

// 			ctx := req.Context()
// 			ctx = dbtrxmiddlewarepkg.NewContext(ctx, dbTrxCtxValue)
// 			req = req.WithContext(ctx)

// 			resprec := httptest.NewRecorder()

// 			router := mux.NewRouter()

// 			router.Name(route.Name).
// 				Methods(route.Method).
// 				Path(route.Path).
// 				HandlerFunc(route.HandlerFunc)

// 			router.ServeHTTP(resprec, req)

// 			if !tc.WantError {
// 				assert.Equal(t, resprec.Code, tc.StatusCode)
// 				returnedToken := tokenhttputilpkg.Token{}
// 				err := json.NewDecoder(resprec.Body).Decode(&returnedToken)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
// 				assert.NotEmpty(t, returnedToken.Text)
// 			} else {
// 				assert.Equal(t, resprec.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestRefreshToken() {
// 	dbTrx := &gorm.DB{}

// 	var authN authpkg.IAuth

// 	userDatastore := datastoremodel.User{}

// 	loginDatastore := datastoremodel.Login{}

// 	authDatastore := datastoremodel.Auth{}

// 	authDetailsCtxValue := domainmodel.Auth{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInRefreshingTheToken",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				userDatastore = datastoremodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&userDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				loginDatastore = datastoremodel.Login{
// 					UserID:   userDatastore.ID,
// 					Username: username,
// 					Password: password,
// 				}

// 				result = dbTrx.Create(&loginDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				authDatastore = datastoremodel.Auth{
// 					UserID: userDatastore.ID,
// 				}

// 				result = dbTrx.Create(&authDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				authDetailsCtxValue = authDatastore.ToDomain()
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsEmpty",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				authDetailsCtxValue = domainmodel.Auth{}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDatastoreRepository := authdatastorerepository.New(dbTrx)

// 			userDatastoreRepository := userdatastorerepository.New(dbTrx)

// 			loginDatastoreRepository := logindatastorerepository.New(dbTrx)

// 			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
// 				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
// 			authHandler := authhandler.New(authService)

// 			route := routehttputilpkg.Route{
// 				Name:        "RefreshToken",
// 				Method:      http.MethodPost,
// 				Path:        "/refresh_token",
// 				HandlerFunc: authHandler.RefreshToken,
// 			}

// 			requestData := requesthttputilpkg.RequestData{
// 				Method: route.Method,
// 				Target: route.Path,
// 			}

// 			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

// 			ctx := req.Context()
// 			ctx = authmiddlewarepkg.NewContext(ctx, authDetailsCtxValue)
// 			req = req.WithContext(ctx)

// 			resprec := httptest.NewRecorder()

// 			router := mux.NewRouter()

// 			router.Name(route.Name).
// 				Methods(route.Method).
// 				Path(route.Path).
// 				HandlerFunc(route.HandlerFunc)

// 			router.ServeHTTP(resprec, req)

// 			if !tc.WantError {
// 				assert.Equal(t, resprec.Code, tc.StatusCode)
// 				returnedToken := tokenhttputilpkg.Token{}
// 				err := json.NewDecoder(resprec.Body).Decode(&returnedToken)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
// 				assert.NotEmpty(t, returnedToken.Text)
// 			} else {
// 				assert.Equal(t, resprec.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestChangePassword() {
// 	dbTrx := &gorm.DB{}

// 	var authN authpkg.IAuth

// 	userDatastore := datastoremodel.User{}

// 	loginDatastore := datastoremodel.Login{}

// 	authDatastore := datastoremodel.Auth{}

// 	auth := domainmodel.Auth{}

// 	passwords := securitypkg.Passwords{}

// 	body := ""

// 	authDetailsCtxValue := domainmodel.Auth{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInResettingThePassword",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				userDatastore = datastoremodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&userDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				loginDatastore = datastoremodel.Login{
// 					UserID:   userDatastore.ID,
// 					Username: username,
// 					Password: password,
// 				}

// 				result = dbTrx.Create(&loginDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				authDatastore = datastoremodel.Auth{
// 					UserID: userDatastore.ID,
// 				}

// 				result = dbTrx.Create(&authDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				auth = authDatastore.ToDomain()

// 				currentPassword := password
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = securitypkg.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"current_password":"%s",
// 					"new_password":"%s"
// 				}`,
// 					passwords.CurrentPassword, passwords.NewPassword)

// 				authDetailsCtxValue = auth
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsEmpty",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				authDetailsCtxValue = domainmodel.Auth{}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				userDatastore = datastoremodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&userDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				loginDatastore = datastoremodel.Login{
// 					UserID:   userDatastore.ID,
// 					Username: username,
// 					Password: password,
// 				}

// 				result = dbTrx.Create(&loginDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				authDatastore = datastoremodel.Auth{
// 					UserID: userDatastore.ID,
// 				}

// 				result = dbTrx.Create(&authDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				auth = authDatastore.ToDomain()

// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = securitypkg.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				body = fmt.Sprintf(`
// 					"current_password":"%s",
// 					"new_password":"%s"
// 				`,
// 					passwords.CurrentPassword, passwords.NewPassword)

// 				authDetailsCtxValue = auth
// 			},
// 			StatusCode: http.StatusBadRequest,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				auth = domainmodel.Auth{}

// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = securitypkg.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"current_password":"%s",
// 					"new_password":"%s"
// 				}`,
// 					passwords.CurrentPassword, passwords.NewPassword)

// 				authDetailsCtxValue = auth
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDatastoreRepository := authdatastorerepository.New(dbTrx)

// 			userDatastoreRepository := userdatastorerepository.New(dbTrx)

// 			loginDatastoreRepository := logindatastorerepository.New(dbTrx)

// 			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
// 				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
// 			authHandler := authhandler.New(authService)

// 			route := routehttputilpkg.Route{
// 				Name:        "ChangePassword",
// 				Method:      http.MethodPost,
// 				Path:        "/change_password",
// 				HandlerFunc: authHandler.ChangePassword,
// 			}

// 			requestData := requesthttputilpkg.RequestData{
// 				Method: route.Method,
// 				Target: route.Path,
// 				Body:   body,
// 			}

// 			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

// 			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

// 			ctx := req.Context()
// 			ctx = authmiddlewarepkg.NewContext(ctx, authDetailsCtxValue)
// 			req = req.WithContext(ctx)

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
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
// 				assert.NotEmpty(t, returnedMessage.Text)
// 			} else {
// 				assert.Equal(t, resprec.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestSignOut() {
// 	dbTrx := &gorm.DB{}

// 	var authN authpkg.IAuth

// 	userDatastore := datastoremodel.User{}

// 	loginDatastore := datastoremodel.Login{}

// 	authDatastore := datastoremodel.Auth{}

// 	auth := domainmodel.Auth{}

// 	authDetailsCtxValue := domainmodel.Auth{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSigningOut",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				userDatastore = datastoremodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&userDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				loginDatastore = datastoremodel.Login{
// 					UserID:   userDatastore.ID,
// 					Username: username,
// 					Password: password,
// 				}

// 				result = dbTrx.Create(&loginDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				authDatastore = datastoremodel.Auth{
// 					UserID: userDatastore.ID,
// 				}

// 				result = dbTrx.Create(&authDatastore)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				auth = authDatastore.ToDomain()

// 				authDetailsCtxValue = auth
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsInvalid",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				authDetailsCtxValue = domainmodel.Auth{}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

// 				authN = authpkg.New(ts.RSAKeys)

// 				auth = domainmodel.Auth{}

// 				authDetailsCtxValue = auth
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDatastoreRepository := authdatastorerepository.New(dbTrx)

// 			userDatastoreRepository := userdatastorerepository.New(dbTrx)

// 			loginDatastoreRepository := logindatastorerepository.New(dbTrx)

// 			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
// 				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
// 			authHandler := authhandler.New(authService)

// 			route := routehttputilpkg.Route{
// 				Name:        "SignOut",
// 				Method:      http.MethodPost,
// 				Path:        "/sign_out",
// 				HandlerFunc: authHandler.SignOut,
// 			}

// 			requestData := requesthttputilpkg.RequestData{
// 				Method: route.Method,
// 				Target: route.Path,
// 			}

// 			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

// 			ctx := req.Context()
// 			ctx = authmiddlewarepkg.NewContext(ctx, authDetailsCtxValue)
// 			req = req.WithContext(ctx)

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
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
// 				assert.NotEmpty(t, returnedMessage.Text)
// 			} else {
// 				assert.Equal(t, resprec.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }
