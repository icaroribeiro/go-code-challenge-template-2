package auth

import (
	"github.com/dgrijalva/jwt-go"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"
)

// IAuth interface is the auth's contract.
type IAuth interface {
	CreateToken(auth domainentity.Auth, tokenExpTimeInSec int) (string, error)
	DecodeToken(tokenString string) (*jwt.Token, error)
	ValidateTokenRenewal(tokenString string, timeBeforeTokenExpTimeInSec int) (*jwt.Token, error)
	FetchAuthFromToken(token *jwt.Token) (domainentity.Auth, error)
}
