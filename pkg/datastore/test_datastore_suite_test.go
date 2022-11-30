package datastore_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	envpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/env"
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

type TestSuite struct {
	suite.Suite
	Cases               Cases
	PostgresDBConfig    map[string]string
	PostgresDBConfigURL string
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

func (ts *TestSuite) SetupSuite() {
	postgresDBDriver := "postgres"
	postresDBUser := envpkg.GetEnvWithDefaultValue("DB_USER", "postgres")
	postresDBPassword := envpkg.GetEnvWithDefaultValue("DB_PASSWORD", "postgres")
	postresDBHost := envpkg.GetEnvWithDefaultValue("DB_HOST", "localhost")
	postresDBPort := envpkg.GetEnvWithDefaultValue("DB_PORT", "5432")
	postresDBName := envpkg.GetEnvWithDefaultValue("DB_NAME", "db")

	ts.PostgresDBConfig = map[string]string{
		"DRIVER":   postgresDBDriver,
		"USER":     postresDBUser,
		"PASSWORD": postresDBPassword,
		"HOST":     postresDBHost,
		"PORT":     postresDBPort,
		"NAME":     postresDBName,
	}

	ts.PostgresDBConfigURL = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		ts.PostgresDBConfig["USER"],
		ts.PostgresDBConfig["PASSWORD"],
		ts.PostgresDBConfig["HOST"],
		ts.PostgresDBConfig["PORT"],
		ts.PostgresDBConfig["NAME"],
	)
}

func TestDatastoreSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
