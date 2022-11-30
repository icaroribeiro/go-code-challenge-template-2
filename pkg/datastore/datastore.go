package datastore

import (
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"gorm.io/gorm"
)

type Provider struct {
	DB *gorm.DB
}

// New is the factory function that encapsulates the implementation related to datastore.
func New(dbConfig map[string]string) (IDatastore, error) {
	driver := dbConfig["DRIVER"]

	switch driver {
	case "postgres":
		return NewPostgresDriver(dbConfig)
	}

	return nil, customerror.Newf("sql database driver %s is not recognized", driver)
}

// GetInstance is the function that gets the database instance.
func (p *Provider) GetInstance() *gorm.DB {
	return p.DB
}

// Close is the function that closes the database connection, releasing any open resources.
func (p *Provider) Close() error {
	db, err := p.DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
