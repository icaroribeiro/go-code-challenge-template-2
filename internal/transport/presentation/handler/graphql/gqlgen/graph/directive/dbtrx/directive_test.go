package dbtrx_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	dbtrxdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/dbtrx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestMiddlewareUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestNewContext() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrxCtxValue := &gorm.DB{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingACopyOfAContextWithAnAssociatedValue",
			SetUp: func(t *testing.T) {
				dbTrxCtxValue = db
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedCtx := dbtrxdirectivepkg.NewContext(ctx, dbTrxCtxValue)

			if !tc.WantError {
				assert.NotEmpty(t, returnedCtx)
				returnedDBTrxCtxValue, ok := dbtrxdirectivepkg.FromContext(returnedCtx)
				assert.True(t, ok, "Unexpected type assertion error")
				assert.Equal(t, dbTrxCtxValue, returnedDBTrxCtxValue)
			}
		})
	}
}

func (ts *TestSuite) TestFromContext() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrxCtxValue := &gorm.DB{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAssociatedValueWithAContext",
			SetUp: func(t *testing.T) {
				dbTrxCtxValue = db
				ctx = dbtrxdirectivepkg.NewContext(ctx, dbTrxCtxValue)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedDBTrxCtxValue, ok := dbtrxdirectivepkg.FromContext(ctx)

			if !tc.WantError {
				assert.True(t, ok, "Unexpected type assertion error")
				assert.NotEmpty(t, returnedDBTrxCtxValue)
				assert.Equal(t, dbTrxCtxValue, returnedDBTrxCtxValue)
			}
		})
	}
}

func (ts *TestSuite) TestDBTrxMiddleware() {
	username := fake.Username()

	user := domainmodel.User{
		Username: username,
	}

	body := fmt.Sprintf(`
	{
		"username":"%s",
	}`,
		username)

	driver := "postgres"
	db, mock := NewMockDB(driver)
	dbAux := &gorm.DB{}

	var handlerFunc func(w http.ResponseWriter, r *http.Request)

	statusCode := 0

	sqlQuery := `INSERT INTO "users" ("id","username","created_at","updated_at") VALUES ($1,$2,$3,$4)`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithDBTrxMiddleware",
			SetUp: func(t *testing.T) {
				dbAux = db

				statusCode = http.StatusOK

				handlerFunc = func(w http.ResponseWriter, r *http.Request) {
					dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

					userDatastore := datastoremodel.User{
						Username: username,
					}

					_ = dbAux.Create(&userDatastore)

					responsehttputilpkg.RespondWithJson(w, http.StatusOK, messagehttputilpkg.Message{Text: "ok"})
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			WantError: false,
		},
		// {
		// 	Context: "ItShouldFailIfTheDatabaseParameterUsedByTheDBTrxMiddlewareIsNil",
		// 	SetUp: func(t *testing.T) {
		// 		dbAux = nil

		// 		statusCode = http.StatusInternalServerError

		// 		handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// 			_, ok := dbtrxmiddlewarepkg.FromContext(r.Context())
		// 			if !ok {
		// 				responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
		// 				return
		// 			}
		// 		}
		// 	},
		// 	WantError: true,
		// },
		// {
		// 	Context: "ItShouldFailIfTheDatabaseTransactionPerformedByTheWrappedFunctionFails",
		// 	SetUp: func(t *testing.T) {
		// 		dbAux = db

		// 		statusCode = http.StatusInternalServerError

		// 		handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// 			dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

		// 			userDatastore := datastoremodel.User{
		// 				Username: username,
		// 			}

		// 			result := dbAux.Create(&userDatastore)
		// 			if result.Error != nil {
		// 				responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
		// 			}
		// 		}

		// 		mock.ExpectBegin()

		// 		mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
		// 			WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
		// 			WillReturnError(customerror.New("failed"))

		// 		mock.ExpectRollback()
		// 	},
		// 	WantError: true,
		// },
		// {
		// 	Context: "ItShouldFailIfTheCommitStatementToEndTheDatabaseTransactionExecutedInsideTheDBTrxMiddlewareFails",
		// 	SetUp: func(t *testing.T) {
		// 		dbAux = db

		// 		statusCode = http.StatusOK

		// 		handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// 			dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

		// 			userDatastore := datastoremodel.User{
		// 				Username: username,
		// 			}

		// 			result := dbAux.Create(&userDatastore)
		// 			if result.Error == nil {
		// 				responsehttputilpkg.RespondWithJson(w, http.StatusOK, messagehttputilpkg.Message{Text: "ok"})
		// 			}
		// 		}

		// 		mock.ExpectBegin()

		// 		mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
		// 			WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
		// 			WillReturnResult(sqlmock.NewResult(1, 1))

		// 		mock.ExpectCommit().WillReturnError(customerror.New("failed"))
		// 	},
		// 	WantError: true,
		// },
		// {
		// 	Context: "ItShouldFailIfTheRollbackStatementToEndTheDatabaseTransactionExecutedInsideTheDBTrxMiddlewareFails",
		// 	SetUp: func(t *testing.T) {
		// 		dbAux = db

		// 		statusCode = http.StatusInternalServerError

		// 		handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// 			dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

		// 			userDatastore := datastoremodel.User{
		// 				Username: username,
		// 			}

		// 			result := dbAux.Create(&userDatastore)
		// 			if result.Error == nil {
		// 				responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
		// 			}
		// 		}

		// 		mock.ExpectBegin()

		// 		mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
		// 			WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
		// 			WillReturnResult(sqlmock.NewResult(1, 1))

		// 		mock.ExpectRollback().WillReturnError(customerror.New("failed"))
		// 	},
		// 	WantError: true,
		// },
		// {
		// 	Context: "ItShouldFailIfTheDatabaseTransactionPerformedByTheWrappedFunctionFailsAndTheFunctionCallsPanicMethodToStopItsExecutionImmediately",
		// 	SetUp: func(t *testing.T) {
		// 		dbAux = db

		// 		statusCode = http.StatusInternalServerError

		// 		handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// 			dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

		// 			userDatastore := datastoremodel.User{
		// 				Username: username,
		// 			}

		// 			result := dbAux.Create(&userDatastore)
		// 			if result.Error != nil {
		// 				// It is duplicated only to test the code that evaluates
		// 				// if the header is already written in the WriteHeader method.
		// 				responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
		// 				responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
		// 			}

		// 			panic(customerror.New("failed"))
		// 		}

		// 		mock.ExpectBegin()

		// 		mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
		// 			WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
		// 			WillReturnError(customerror.New("failed"))

		// 		mock.ExpectRollback()
		// 	},
		// 	WantError: true,
		// },
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			// dbtrxMiddleware := dbtrxmiddlewarepkg.DBTrx(dbAux)

			// returnedHandlerFunc := adapterhttputilpkg.AdaptFunc(handlerFunc).With(dbtrxMiddleware)

			// route := routehttputilpkg.Route{
			// 	Name:        "Testing",
			// 	Method:      http.MethodGet,
			// 	Path:        "/testing",
			// 	HandlerFunc: returnedHandlerFunc,
			// }

			// requestData := requesthttputilpkg.RequestData{
			// 	Method: route.Method,
			// 	Target: route.Path,
			// 	Body:   body,
			// }

			// reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

			// req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

			// resprec := httptest.NewRecorder()

			// router := mux.NewRouter()

			// router.Name(route.Name).
			// 	Methods(route.Method).
			// 	Path(route.Path).
			// 	HandlerFunc(route.HandlerFunc)

			// router.ServeHTTP(resprec, req)

			if !tc.WantError {
				// assert.Equal(t, resprec.Result().Header.Get("Content-Type"), "application/json")
				// assert.Equal(t, statusCode, resprec.Result().StatusCode)
				// returnedMessage := messagehttputilpkg.Message{}
				// err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				// assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				// assert.NotEmpty(t, returnedMessage.Text)
			} else {
				// assert.Equal(t, statusCode, resprec.Result().StatusCode)
			}

			// err := mock.ExpectationsWereMet()
			// assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}
