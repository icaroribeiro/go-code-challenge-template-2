package model

import (
	"time"

	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/security"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login is the representation of the login's datastore model.
type Login struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" validate:"uuid"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;unique" validate:"uuid"`
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsEmpty is the function that checks if login's datastore model is empty.
func (l Login) IsEmpty() bool {
	return l == Login{}
}

// BeforeCreate is a Gorm hook that is called before creating a login in the datastore.
func (l *Login) BeforeCreate(tx *gorm.DB) error {
	l.ID = uuid.NewV4()

	security := securitypkg.New()

	hashedPassword, err := security.HashPassword(l.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	l.Password = string(hashedPassword)

	return nil
}

// BeforeUpdate is a Gorm hook that is called before updating a login in the datastore.
func (l *Login) BeforeUpdate(tx *gorm.DB) error {
	security := securitypkg.New()

	hashedPassword, err := security.HashPassword(l.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	l.Password = string(hashedPassword)

	return nil
}

// FromDomain is the function that builds a datastore model based on the model's data from domain.
func (l *Login) FromDomain(login domainmodel.Login) {
	l.ID = login.ID
	l.UserID = login.UserID
	l.Username = login.Username
	l.Password = login.Password
}

// ToDomain is the function that returns a domain model built using the model's data from datastore.
func (l *Login) ToDomain() domainmodel.Login {
	return domainmodel.Login{
		ID:       l.ID,
		UserID:   l.UserID,
		Username: l.Username,
		Password: l.Password,
	}
}
