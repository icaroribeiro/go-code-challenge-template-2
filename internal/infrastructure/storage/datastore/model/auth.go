package model

import (
	"time"

	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Auth is the representation of the auth's datastore model.
type Auth struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;unique"`
	CreatedAt time.Time
}

// BeforeCreate is a Gorm hook that is called before create an auth in the database.
func (a *Auth) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.NewV4()

	return nil
}

// FromDomain is the function that builds a database model based on the model's data from domain.
func (a *Auth) FromDomain(auth domainmodel.Auth) {
	a.ID = auth.ID
	a.UserID = auth.UserID
}

// ToDomain is the function that returns a domain model built using the model's data from database.
func (a *Auth) ToDomain() domainmodel.Auth {
	return domainmodel.Auth{
		ID:     a.ID,
		UserID: a.UserID,
	}
}
