package datastore

import "gorm.io/gorm"

// IDatastore interface is the datastore's contract.
type IDatastore interface {
	GetInstance() *gorm.DB
	Close() error
}
