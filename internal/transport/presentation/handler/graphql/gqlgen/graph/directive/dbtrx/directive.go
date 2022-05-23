package dbtrx

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"gorm.io/gorm"
)

var dbTrxCtxKey = &contextKey{"db_trx"}
var dbTrxStateCtxKey = &contextKey{"db_trx_state"}

type contextKey struct {
	name string
}

type DBTrxState struct {
	DBTrx      *gorm.DB
	NeedCommit bool
}

// IsEmpty is the function that checks if dbTrxState's model is empty.
func (d DBTrxState) IsEmpty() bool {
	return d == DBTrxState{}
}

// NewContext is the function that returns a new Context that carries db_trx_state value.
func NewContext(ctx context.Context, dbTrx *gorm.DB) context.Context {
	return context.WithValue(ctx, dbTrxCtxKey, dbTrx)
}

// FromContext is the function that returns the db_trx_state value stored in context, if any.
func FromContext(ctx context.Context) (*gorm.DB, bool) {
	log.Println(ctx)
	raw, ok := ctx.Value(dbTrxCtxKey).(*gorm.DB)
	log.Println(ok)
	return raw, ok
}

// UseDBTrx is the function that...
func UseDBTrx(db *gorm.DB) func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		if db == nil {
			return nil, customerror.New("Database is nil")
		}

		dbTrx := db.Begin()
		// It is necessary to set a struct with database transaction hat can be used for performing operations with transaction.
		// dbTrxState := DBTrxState{
		// 	DBTrx:      dbTrx,
		// 	NeedCommit: false,
		// }

		ctx = NewContext(ctx, dbTrx)

		log.Println("UseDBTrx")

		return next(ctx)
	}
}
