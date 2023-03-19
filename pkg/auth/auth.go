package auth

import (
	"crypto/rsa"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	uuid "github.com/satori/go.uuid"
)

type Auth struct {
	RSAKeys RSAKeys
}

// RSAKeys is the representation of the RSA keys.
type RSAKeys struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func New(rsaKeys RSAKeys) IAuth {
	return &Auth{
		RSAKeys: RSAKeys{
			PublicKey:  rsaKeys.PublicKey,
			PrivateKey: rsaKeys.PrivateKey,
		},
	}
}

// CreateToken is the function that creates a new token for a specific auth and time duration.
func (a *Auth) CreateToken(auth domainentity.Auth, tokenExpTimeInSec int) (string, error) {
	duration := time.Second * time.Duration(tokenExpTimeInSec)

	claims := jwt.MapClaims{
		"auth_id":    auth.ID.String(),
		"user_id":    auth.UserID.String(),
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(duration).Unix(),
		"authorized": true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(a.RSAKeys.PrivateKey)
}

func parseToken(tokenString string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, customerror.New("unexpected signing method when trying to decode the token")
		}
		return publicKey, nil
	}

	return jwt.Parse(tokenString, keyFunc)
}

// ExtractTokenString is the function that extracts the token string from authentication header string.
func (a *Auth) ExtractTokenString(authHeaderString string) (string, error) {
	if len(authHeaderString) == 0 {
		errorMessage := "the auth header must be informed along with the token"
		return "", customerror.BadRequest.New(errorMessage)
	}

	bearerToken := strings.Split(authHeaderString, " ")
	if bearerToken[1] == "" {
		errorMessage := "the token must be associated with the auth header"
		return "", customerror.BadRequest.New(errorMessage)
	}

	return bearerToken[1], nil
}

// DecodeToken is the function that translates a token string in a jwt token
// and checks if the jwt token is valid or not.
func (a *Auth) DecodeToken(tokenString string) (*jwt.Token, error) {
	token, err := parseToken(tokenString, a.RSAKeys.PublicKey)

	if verr, ok := err.(*jwt.ValidationError); ok {
		switch verr.Errors {
		case jwt.ValidationErrorExpired:
			errorMessage := "the token has expired"
			return nil, customerror.Unauthorized.New(errorMessage)
		default:
			return nil, err
		}
	}

	return token, nil
}

// ValidateTokenRenewal is the function that validates if the jwt token is already expired to be renewed.
func (a *Auth) ValidateTokenRenewal(token *jwt.Token, timeBeforeTokenExpTimeInSec int) (*jwt.Token, error) {
	claims, _ := token.Claims.(jwt.MapClaims)

	expiredAt, _ := claims["exp"].(int64)

	duration := time.Second * time.Duration(timeBeforeTokenExpTimeInSec)

	if time.Until(time.Unix(int64(expiredAt), 0)) > duration {
		errorMessage := "the token expiration time is not within the time prior to the time before token expiration time"
		return token, customerror.BadRequest.New(errorMessage)
	}

	return token, nil
}

// FetchAuthFromToken is the function that get auth data from the token.
func (a *Auth) FetchAuthFromToken(token *jwt.Token) (domainentity.Auth, error) {
	auth := domainentity.Auth{}

	if token == nil {
		return auth, customerror.New("the token is nil")
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	id, ok := claims["auth_id"].(string)
	if !ok {
		return auth, customerror.New("failed to extract the auth_id from the token")
	}

	authID, err := uuid.FromString(id)
	if err != nil {
		return auth, customerror.Newf("failed to convert the auth_id %s from the token to UUID", id)
	}

	id, ok = claims["user_id"].(string)
	if !ok {
		return auth, customerror.New("failed to extract the user_id from the token")
	}

	userID, err := uuid.FromString(id)
	if err != nil {
		return auth, customerror.Newf("failed to convert the user_id %s from the to UUID", id)
	}

	auth.ID = authID
	auth.UserID = userID

	return auth, nil
}
