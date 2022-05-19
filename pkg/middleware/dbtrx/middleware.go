package dbtrx

import (
	"context"
	"log"
	"net/http"

	"gorm.io/gorm"
)

// A private key for context that only this package can access.
// This is important to prevent collisions between different context uses.
var dbTrxCtxKey = &contextKey{"db_trx"}

type contextKey struct {
	name string
}

// DBTrx is the function that  wraps a http.Handler to enable using a database transaction during an API incoming request.
func DBTrx(db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if db == nil {
				next.ServeHTTP(w, r)
				return
			}

			dbTrx := db.Begin()
			defer func() {
				if r := recover(); r != nil {
					dbTrx.Rollback()
				}
			}()

			// It is necessary to set database transaction that can be used for performing operations with transaction.
			ctx := context.WithValue(r.Context(), dbTrxCtxKey, dbTrx)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)

			// If graphQL not return error message {
			if err := dbTrx.Commit().Error; err != nil {
				errorMsg := "database transaction commit failed"
				log.Panicf("%s: %s", errorMsg, err.Error())
			}
			// } else {
			// 	log.Printf("rolling back database transaction due to status code: %d", wrapped.statusCode)
			// 	dbTrx.Rollback()
			// }
		}
	}
}

// ForContext is the function that finds the db_trx from the context.
func ForContext(ctx context.Context) *gorm.DB {
	raw, _ := ctx.Value(dbTrxCtxKey).(*gorm.DB)
	return raw
}
