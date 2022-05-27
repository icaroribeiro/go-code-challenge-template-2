package resolver_test

import (
	"context"
	"log"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/DATA-DOG/go-sqlmock"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	authdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/auth"
	dbtrxdirectivepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql/gqlgen/graph/directive/dbtrx"
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

func MockDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		return next(ctx)
	}
}

func AddAuthDetailsToCtx(ctx context.Context, auth domainmodel.Auth) client.Option {
	return func(bd *client.Request) {
		ctx := authdirectivepkg.NewContext(ctx, auth)
		bd.HTTP = bd.HTTP.WithContext(ctx)
	}
}

func AddDBTrxToCtx(ctx context.Context, dbTrx *gorm.DB) client.Option {
	return func(bd *client.Request) {
		ctx := dbtrxdirectivepkg.NewContext(ctx, dbTrx)
		bd.HTTP = bd.HTTP.WithContext(ctx)
	}
}
