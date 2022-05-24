package dbtrx

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"gorm.io/gorm"
)

var dbTrxCtxKey = &contextKey{"db_trx"}

type contextKey struct {
	name string
}

// NewContext is the function that returns a new Context that carries db_trx_state value.
func NewContext(ctx context.Context, dbTrx *gorm.DB) context.Context {
	return context.WithValue(ctx, dbTrxCtxKey, dbTrx)
}

// FromContext is the function that returns the db_trx_state value stored in context, if any.
func FromContext(ctx context.Context) (*gorm.DB, bool) {
	raw, ok := ctx.Value(dbTrxCtxKey).(*gorm.DB)
	return raw, ok
}

// DBTrxMiddleware the function that acts as a HTTP middleware to enable using a database transaction during an API incoming request.
func DBTrxMiddleware(db *gorm.DB) func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		if db == nil {
			return next(ctx)
		}

		dbTrx := db.Begin()

		defer func() {
			if r := recover(); r != nil {
				var err error
				switch r := r.(type) {
				case error:
					err = r
				default:
					err = customerror.Newf("%v", r)
				}
				log.Printf("Transaction is being rolled back: %s \n", err.Error())
				dbTrx.Rollback()
				return
			}
		}()

		// It is necessary to set database transaction that can be used for performing operations with transaction.
		ctx = NewContext(ctx, dbTrx)

		res, err := next(ctx)

		if err == nil {
			if err = dbTrx.Commit().Error; err != nil {
				log.Printf("failed to commit database transaction: %s", err.Error())
			}
		} else {
			log.Printf("database transaction is being rolled back: %s", err.Error())
			if err := dbTrx.Rollback().Error; err != nil {
				log.Printf("failed to rollback database transaction: %s", err.Error())
			}
		}

		return res, err
	}
}
