package datastore

import (
	"fmt"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDriver struct {
	Provider Provider
}

// NewPostgresDriver is the factory function that encapsulates the implementation related to postgres.
func NewPostgresDriver(dbConfig map[string]string) (IDatastore, error) {
	dsn := ""

	if dbConfig["URL"] != "" {
		dsn = dbConfig["URL"]
	} else {
		dsn = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			dbConfig["USER"],
			dbConfig["PASSWORD"],
			dbConfig["HOST"],
			dbConfig["PORT"],
			dbConfig["NAME"],
		)
	}

	dialector := postgres.Open(dsn)

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return &PostgresDriver{}, customerror.Newf("failed to establish a database connection: %s", err.Error())
	}

	return &PostgresDriver{
		Provider{
			DB: db,
		},
	}, nil
}

// GetInstance is the function that gets the database instance.
func (d *PostgresDriver) GetInstance() *gorm.DB {
	return d.Provider.GetInstance()
}

// Close is the function that closes the database connection, releasing any open resources.
func (d *PostgresDriver) Close() error {
	return d.Provider.Close()
}
