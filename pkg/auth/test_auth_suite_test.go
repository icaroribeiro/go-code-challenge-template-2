package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	"github.com/stretchr/testify/suite"
)

type Case struct {
	Context   string
	SetUp     func(t *testing.T)
	WantError bool
	TearDown  func(t *testing.T)
}

type Cases []Case

type ReturnArgs [][]interface{}

type TestSuite struct {
	suite.Suite
	RSAKeys authpkg.RSAKeys
	Cases   Cases
}

func (ts *TestSuite) SetupSuite() {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, fake.RandomInt([]int{1024, 2048}))
	if err != nil {
		log.Panicf("failed to parse the RSA public key: %s", err.Error())
	}

	rsaPublicKey := &rsaPrivateKey.PublicKey

	ts.RSAKeys = authpkg.RSAKeys{
		PublicKey:  rsaPublicKey,
		PrivateKey: rsaPrivateKey,
	}
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
