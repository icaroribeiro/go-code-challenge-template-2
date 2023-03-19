package dbtrx_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/DATA-DOG/go-sqlmock"
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	persistententity "github.com/icaroribeiro/go-code-challenge-template-2/internal/infrastructure/datastore/perentity"
	dbtrxdirective "github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/directive/dbtrx"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestDBTrxMiddleware() {
	user := domainentity.UserFactory(nil)

	driver := "postgres"
	db, mock := NewMockDB(driver)
	dbAux := &gorm.DB{}

	ctx := context.Background()

	var next graphql.Resolver

	sqlQuery := `INSERT INTO "users" ("username","created_at","updated_at","id") VALUES ($1,$2,$3,$4)`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithDBTrxMiddleware",
			SetUp: func(t *testing.T) {
				dbAux = db

				next = func(ctx context.Context) (interface{}, error) {
					dbAux, _ := dbtrxdirective.FromContext(ctx)

					userDatastore := persistententity.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)

					return nil, result.Error
				}

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(user.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

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

					userDatastore := persistententity.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						return nil, result.Error
					}

					return nil, nil
				}

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(user.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

					userDatastore := persistententity.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						return nil, result.Error
					}

					return nil, nil
				}

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(user.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

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

					userDatastore := persistententity.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						return nil, result.Error
					}

					return nil, nil
				}

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(user.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

					userDatastore := persistententity.User{
						Username: user.Username,
					}

					_ = dbAux.Create(&userDatastore)

					panic(customerror.New("failed"))
				}

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(user.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

					userDatastore := persistententity.User{
						Username: user.Username,
					}

					_ = dbAux.Create(&userDatastore)

					panic("failed")
				}

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(user.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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
