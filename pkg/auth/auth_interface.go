package auth

import (
	"github.com/dgrijalva/jwt-go"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
)

// IAuth interface is the auth's contract.
type IAuth interface {
	CreateToken(auth domainmodel.Auth, tokenExpTimeInSec int) (string, error)
	DecodeToken(tokenString string) (*jwt.Token, error)
	ValidateTokenRenewal(tokenString string, timeBeforeTokenExpTimeInSec int) (*jwt.Token, error)
	FetchAuthFromToken(token *jwt.Token) (domainmodel.Auth, error)
}
