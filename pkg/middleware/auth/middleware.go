package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
	"gorm.io/gorm"
)

var authDetailsCtxKey = &contextKey{"auth_details"}

type contextKey struct {
	name string
}

// NewContext is the function that returns a new Context that carries auth_details value.
func NewContext(ctx context.Context, auth domainentity.Auth) context.Context {
	return context.WithValue(ctx, authDetailsCtxKey, auth)
}

// FromContext is the function that returns the auth_details value stored in context, if any.
func FromContext(ctx context.Context) (domainentity.Auth, bool) {
	raw, ok := ctx.Value(authDetailsCtxKey).(domainentity.Auth)
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

func buildAuth(db *gorm.DB, authN authpkg.IAuth, token *jwt.Token) (domainentity.Auth, error) {
	auth, err := authN.FetchAuthFromToken(token)
	if err != nil {
		return domainentity.Auth{}, err
	}

	// Before proceeding is necessary to check if the user who is performing operations is logged
	// based on the authentication details inserted within in the token.
	authAux := domainentity.Auth{}

	result := db.Find(&authAux, "id=?", auth.ID)
	if result.Error != nil {
		return domainentity.Auth{}, result.Error
	}

	if authAux.IsEmpty() {
		errorMessage := "you are not logged in, then perform a login to get a token before proceeding"
		return domainentity.Auth{}, customerror.BadRequest.New(errorMessage)
	}

	if auth.UserID.String() != authAux.UserID.String() {
		errorMessage := "the token's auth_id and user_id are not associated"
		return domainentity.Auth{}, customerror.BadRequest.New(errorMessage)
	}

	return auth, nil
}

// Auth is the function that wraps a http.Handler to evaluate the authentication of API based on a JWT token.
func Auth(db *gorm.DB, authN authpkg.IAuth) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := extractTokenString(w, r)
			if err != nil {
				responsehttputilpkg.RespondErrorWithJson(w, err)
				return
			}

			token, err := authN.DecodeToken(tokenString)
			if err != nil {
				responsehttputilpkg.RespondErrorWithJson(w, customerror.Unauthorized.New(err.Error()))
				return
			}

			auth, err := buildAuth(db, authN, token)
			if err != nil {
				responsehttputilpkg.RespondErrorWithJson(w, err)
				return
			}

			// It is necessary to set auth details that can be used for performing authenticated operations.
			ctx := NewContext(r.Context(), auth)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	}
}

// AuthRenewal is the function that wraps a http.Handler to evaluate the authentication renewal of API based on a JWT token.
func AuthRenewal(db *gorm.DB, authN authpkg.IAuth, timeBeforeTokenExpTimeInSec int) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := extractTokenString(w, r)
			if err != nil {
				responsehttputilpkg.RespondErrorWithJson(w, err)
				return
			}

			token, err := authN.ValidateTokenRenewal(tokenString, timeBeforeTokenExpTimeInSec)
			if err != nil {
				responsehttputilpkg.RespondErrorWithJson(w, customerror.Unauthorized.New(err.Error()))
				return
			}

			auth, err := buildAuth(db, authN, token)
			if err != nil {
				responsehttputilpkg.RespondErrorWithJson(w, err)
				return
			}

			// It is necessary to set auth details that can be used for performing authenticated operations.
			ctx := NewContext(r.Context(), auth)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	}
}
