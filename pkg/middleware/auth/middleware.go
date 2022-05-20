package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/response"
	"gorm.io/gorm"
)

// A private key for context that only this package can access.
// This is important to prevent collisions between different context uses.
var authDetailsCtxKey = &contextKey{"auth_details"}

type contextKey struct {
	name string
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

func buildAuth(db *gorm.DB, authN authpkg.IAuth, token *jwt.Token) (domainmodel.Auth, error) {
	auth, err := authN.FetchAuthFromToken(token)
	if err != nil {
		return domainmodel.Auth{}, err
	}

	// Before proceeding is necessary to check if the user who is performing operations is logged
	// based on the authentication details inserted within in the token.
	authAux := domainmodel.Auth{}

	result := db.Find(&authAux, "id=?", auth.ID)
	if result.Error != nil {
		return domainmodel.Auth{}, result.Error
	}

	if authAux.IsEmpty() {
		errorMessage := "you are not logged in, then perform a login to get a token before proceeding"
		return domainmodel.Auth{}, customerror.BadRequest.New(errorMessage)
	}

	if auth.UserID.String() != authAux.UserID.String() {
		errorMessage := "the token's auth_id and user_id are not associated"
		return domainmodel.Auth{}, customerror.BadRequest.New(errorMessage)
	}

	return auth, nil
}

func setupAuthDetailsInRequest(w http.ResponseWriter, r *http.Request, auth domainmodel.Auth) *http.Request {
	// It is necessary to set auth details that can be used for performing authenticated operations.
	ctx := context.WithValue(r.Context(), authDetailsCtxKey, auth)
	return r.WithContext(ctx)
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

			r = setupAuthDetailsInRequest(w, r, auth)

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

			r = setupAuthDetailsInRequest(w, r, auth)

			next.ServeHTTP(w, r)
		}
	}
}
