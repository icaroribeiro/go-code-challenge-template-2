package auth

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/response"
	"gorm.io/gorm"
)

var tokenCtxKey = &contextKey{"token"}

type contextKey struct {
	name string
}

// Auth is the function that wraps a http.Handler to evaluate the authentication of API based on a JWT token.
func Auth(db *gorm.DB, authN authpkg.IAuth) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeaderString := r.Header.Get("Authorization")

			tokenString, _ := authN.ExtractTokenString(authHeaderString)
			// if err != nil {
			// 	responsehttputilpkg.RespondErrorWithJson(w, err)
			// 	return
			// }

			if tokenString != "" {
				token, err := authN.DecodeToken(tokenString)
				if err != nil {
					responsehttputilpkg.RespondErrorWithJson(w, customerror.Unauthorized.New(err.Error()))
					return
				}

				// It is necessary to set token that can be used for performing authenticated operations.
				ctx := NewContext(r.Context(), token)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		}
	}
}

// NewContext is the function that returns a new Context that carries auth_details value.
func NewContext(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, tokenCtxKey, token)
}

// FromContext is the function that returns the token value stored in context, if any.
func FromContext(ctx context.Context) (*jwt.Token, bool) {
	raw, ok := ctx.Value(tokenCtxKey).(*jwt.Token)
	return raw, ok
}
