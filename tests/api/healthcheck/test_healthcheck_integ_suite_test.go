package healthcheck_test

import (
	"context"
	"log"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	datastorepkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/datastore"
	envpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/env"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Case struct {
	Context   string
	SetUp     func(t *testing.T)
	WantError bool
	TearDown  func(t *testing.T)
}

type Cases []Case

type TestSuite struct {
	suite.Suite
	DB    *gorm.DB
	Cases Cases
}

type GetHealthCheckQueryResponse struct {
	GetHealthCheck struct {
		Status string
	}
}

var getHealthCheckQuery = `query {
	getHealthCheck { 
		status
	}
}`

var (
	dbDriver   = envpkg.GetEnvWithDefaultValue("DB_DRIVER", "postgres")
	dbUser     = envpkg.GetEnvWithDefaultValue("DB_USER", "postgres")
	dbPassword = envpkg.GetEnvWithDefaultValue("DB_PASSWORD", "postgres")
	dbHost     = envpkg.GetEnvWithDefaultValue("DB_HOST", "localhost")
	dbPort     = envpkg.GetEnvWithDefaultValue("DB_PORT", "5434")
	dbName     = envpkg.GetEnvWithDefaultValue("DB_NAME", "testdb")
)

func setupDBConfig() (map[string]string, error) {
	dbConfig := map[string]string{
		"DRIVER":   dbDriver,
		"USER":     dbUser,
		"PASSWORD": dbPassword,
		"HOST":     dbHost,
		"PORT":     dbPort,
		"NAME":     dbName,
	}

	return dbConfig, nil
}

func MockDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		return next(ctx)
	}
}

func (ts *TestSuite) SetupSuite() {
	dbConfig, err := setupDBConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	datastore, err := datastorepkg.New(dbConfig)
	if err != nil {
		log.Panic(err.Error())
	}

	ts.DB = datastore.GetInstance()
	if ts.DB == nil {
		log.Panicf("The database instance is null")
	}

	if err = ts.DB.Error; err != nil {
		log.Panicf("Got error when acessing the database instance: %s", err.Error())
	}
}

func (ts *TestSuite) TearDownSuite() {
	db, err := ts.DB.DB()
	if err != nil {
		log.Panicf("Got error when acessing *sql.DB from database instance: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		log.Panicf("Got error when closing the database instance: %s", err.Error())
	}
}

func TestHealthCheckIntegSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
