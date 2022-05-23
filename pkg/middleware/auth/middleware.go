package auth

import (
	"context"
	"net/http"
	"strings"

	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	"gorm.io/gorm"
)

var tokenStringCtxKey = &contextKey{"token_string"}

type contextKey struct {
	name string
}

// NewContext is the function that returns a new Context that carries token_string value.
func NewContext(ctx context.Context, tokenString string) context.Context {
	return context.WithValue(ctx, tokenStringCtxKey, tokenString)
}

// FromContext is the function that returns the token_string value stored in context, if any.
func FromContext(ctx context.Context) (string, bool) {
	raw, ok := ctx.Value(tokenStringCtxKey).(string)
	return raw, ok
}

func extractTokenString(w http.ResponseWriter, r *http.Request) (string, error) {
	hdrAuth := r.Header.Get("Authorization")
	if len(hdrAuth) == 0 {
		errorMessage := "the auth header must be informed along with the token"
		return "", customerror.BadRequest.New(errorMessage)
	}

	bearerToken := strings.Split(hdrAuth, " ")
	if len(bearerToken) != 2 {
		errorMessage := "the token must be associated with the auth header"
		return "", customerror.BadRequest.New(errorMessage)
	}

	return bearerToken[1], nil
}

// Auth is the function that wraps a http.Handler to evaluate the authentication of API based on a JWT token.
func Auth(db *gorm.DB, authN authpkg.IAuth) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := extractTokenString(w, r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := NewContext(r.Context(), tokenString)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	}
}
