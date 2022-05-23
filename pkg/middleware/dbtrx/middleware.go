package dbtrx

// import (
// 	"context"
// 	"net/http"

// 	"gorm.io/gorm"
// )

// var dbTrxCtxKey = &contextKey{"db"}

// type contextKey struct {
// 	name string
// }

// // FromContext is the function that returns the db_trx value stored in context, if any.
// func FromContext(ctx context.Context) (*gorm.DB, bool) {
// 	raw, ok := ctx.Value(dbTrxCtxKey).(*gorm.DB)
// 	return raw, ok
// }

// // DB is the function that  wraps a http.Handler to enable using a database instance during an API incoming request.
// func DB(dbInstance *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
// 	return func(next http.HandlerFunc) http.HandlerFunc {
// 		return func(w http.ResponseWriter, r *http.Request) {
// 			if dbInstance == nil {
// 				next.ServeHTTP(w, r)
// 				return
// 			}

// 			ctx := context.WithValue(r.Context(), dbTrxCtxKey, dbTrx)
// 			r = r.WithContext(ctx)

// 			next.ServeHTTP(w, r)
// 		}
// 	}
// }
