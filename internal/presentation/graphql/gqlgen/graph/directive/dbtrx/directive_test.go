package dbtrx_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/DATA-DOG/go-sqlmock"
	datastoreentity "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/entity"
	dbtrxdirective "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/presentation/graphql/gqlgen/graph/directive/dbtrx"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	domainentityfactory "github.com/icaroribeiro/new-go-code-challenge-template-2/tests/factory/core/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestDirectiveUnit(t *testing.T) {
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

			returnedCtx := dbtrxdirective.NewContext(ctx, dbTrxCtxValue)

			if !tc.WantError {
				assert.NotEmpty(t, returnedCtx)
				returnedDBTrxCtxValue, ok := dbtrxdirective.FromContext(returnedCtx)
				assert.True(t, ok, "Unexpected type assertion error.")
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
			Context: "ItShouldSucceedInGettingAnAssociatedValueFromAContext",
			SetUp: func(t *testing.T) {
				dbTrxCtxValue = db
				ctx = dbtrxdirective.NewContext(ctx, dbTrxCtxValue)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedDBTrxCtxValue, ok := dbtrxdirective.FromContext(ctx)

			if !tc.WantError {
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.NotEmpty(t, returnedDBTrxCtxValue)
				assert.Equal(t, dbTrxCtxValue, returnedDBTrxCtxValue)
			}
		})
	}
}

func (ts *TestSuite) TestDBTrxMiddleware() {
	user := domainentityfactory.NewUser(nil)

	driver := "postgres"
	db, mock := NewMockDB(driver)
	dbAux := &gorm.DB{}

	ctx := context.Background()

	var next graphql.Resolver

	sqlQuery := `INSERT INTO "users" ("id","username","created_at","updated_at") VALUES ($1,$2,$3,$4)`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithDBTrxMiddleware",
			SetUp: func(t *testing.T) {
				dbAux = db

				next = func(ctx context.Context) (interface{}, error) {
					dbAux, _ := dbtrxdirective.FromContext(ctx)

					userDatastore := datastoreentity.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)

					return nil, result.Error
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheDatabaseParameterUsedByTheDBTrxMiddlewareIsNil",
			SetUp: func(t *testing.T) {
				dbAux = nil

				next = func(ctx context.Context) (interface{}, error) {
					_, ok := dbtrxdirective.FromContext(ctx)
					if !ok {
						return nil, customerror.New("failed")
					}

					return nil, nil
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionPerformedByTheWrappedFunctionFails",
			SetUp: func(t *testing.T) {
				dbAux = db

				next = func(ctx context.Context) (interface{}, error) {
					dbAux, _ := dbtrxdirective.FromContext(ctx)

					userDatastore := datastoreentity.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						return nil, result.Error
					}

					return nil, nil
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheCommitStatementToEndTheDatabaseTransactionExecutedInsideTheDBTrxMiddlewareFails",
			SetUp: func(t *testing.T) {
				dbAux = db

				next = func(ctx context.Context) (interface{}, error) {
					dbAux, _ := dbtrxdirective.FromContext(ctx)

					userDatastore := datastoreentity.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						return nil, result.Error
					}

					return nil, nil
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(customerror.New("failed"))
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheRollbackStatementToEndTheDatabaseTransactionExecutedInsideTheDBTrxMiddlewareFails",
			SetUp: func(t *testing.T) {
				dbAux = db

				next = func(ctx context.Context) (interface{}, error) {
					dbAux, _ := dbtrxdirective.FromContext(ctx)

					userDatastore := datastoreentity.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						return nil, result.Error
					}

					return nil, nil
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback().WillReturnError(customerror.New("failed"))
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionPerformedByTheWrappedFunctionFailsAndTheFunctionCallsPanicMethodWithErrorParameterToStopItsExecutionImmediately",
			SetUp: func(t *testing.T) {
				dbAux = db

				next = func(ctx context.Context) (interface{}, error) {
					dbAux, _ := dbtrxdirective.FromContext(ctx)

					userDatastore := datastoreentity.User{
						Username: user.Username,
					}

					_ = dbAux.Create(&userDatastore)

					panic(customerror.New("failed"))
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()
			},
			WantError:   true,
			ShouldPanic: true,
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionPerformedByTheWrappedFunctionFailsAndTheFunctionCallsPanicMethodWithNonErrorParameterToStopItsExecutionImmediately",
			SetUp: func(t *testing.T) {
				dbAux = db

				next = func(ctx context.Context) (interface{}, error) {
					dbAux, _ := dbtrxdirective.FromContext(ctx)

					userDatastore := datastoreentity.User{
						Username: user.Username,
					}

					_ = dbAux.Create(&userDatastore)

					panic("failed")
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()
			},
			WantError:   true,
			ShouldPanic: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			dbTrxDirective := dbtrxdirective.New(dbAux)

			_, err := dbTrxDirective.DBTrxMiddleware()(ctx, nil, next)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				if tc.ShouldPanic {
					ShouldPanic(t, next, ctx)
				} else {
					assert.NotNil(t, err, "Predicted error lost.")
				}
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}

func ShouldPanic(t *testing.T, f graphql.Resolver, ctx context.Context) {
	defer func() { recover() }()
	f(ctx)
	t.Errorf("It should have panicked.")
}
