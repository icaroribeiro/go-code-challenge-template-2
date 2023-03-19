package perentity

import (
	"time"

	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// User is the representation of the user's datastore model.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Users is a slice of User.
type Users []User

// BeforeCreate is a Gorm hook that is called before create a user in the datastore.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.NewV4()

	return nil
}

// FromDomain is the function that builds a datastore model based on the model's data from domain.
func (u *User) FromDomain(user domainentity.User) {
	u.ID = user.ID
	u.Username = user.Username
}

// FromDomain is the function that builds a datastore model slice based on the model slice's data from domain.
func (us *Users) FromDomain(users domainentity.Users) {
	u := User{}

	for _, user := range users {
		u.FromDomain(user)
		*us = append(*us, u)
	}
}

// ToDomain is the function that returns a domain model built using the model's data from datastore.
func (u *User) ToDomain() domainentity.User {
	return domainentity.User{
		ID:       u.ID,
		Username: u.Username,
	}
}

// ToDomain is the function that returns a slice of domain model built using the model slice's data from datastore.
func (us *Users) ToDomain() domainentity.Users {
	users := domainentity.Users{}

	for _, u := range *us {
		users = append(users, u.ToDomain())
	}

	return users
}
