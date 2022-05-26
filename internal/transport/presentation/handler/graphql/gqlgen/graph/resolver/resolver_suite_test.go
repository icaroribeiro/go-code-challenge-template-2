package resolver_test

import (
	"context"
	"log"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/DATA-DOG/go-sqlmock"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/auth"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	Cases Cases
}

func addTokenStringCtxValue(ctx context.Context, tokenString string) client.Option {
	return func(bd *client.Request) {
		ctx := authmiddlewarepkg.NewContext(ctx, tokenString)
		bd.HTTP = bd.HTTP.WithContext(ctx)
	}
}

func NewMockDB(driver string) (*gorm.DB, sqlmock.Sqlmock) {
	errorMsg := "failed to open a stub database connection"

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		log.Panicf("%s: %s", errorMsg, err.Error())
	}

	if sqlDB == nil {
		log.Panicf("%s: the sqlDB is null", errorMsg)
	}

	if mock == nil {
		log.Panicf("%s: the mock is null", errorMsg)
	}

	errorMsg = "failed to initialize db session"

	var db *gorm.DB

	switch driver {
	case "postgres":
		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		if err != nil {
			log.Panicf("%s: %s", errorMsg, err.Error())
		}
	}

	if db == nil {
		log.Panicf("%s: the database is null", errorMsg)
	}

	if err = db.Error; err != nil {
		log.Panicf("%s: %s", errorMsg, err.Error())
	}

	return db, mock
}
