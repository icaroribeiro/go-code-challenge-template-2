package presentity

import (
	domainentity "github.com/icaroribeiro/go-code-challenge-template-2/internal/core/domain/entity"
	uuid "github.com/satori/go.uuid"
)

// User is the representation of user's http entity.
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// Users is a slice of User.
type Users []User

// type Users []*User

// FromDomain is the function that builds a http model based on the model's data from domain.
func (u *User) FromDomain(user domainentity.User) {
	u.ID = user.ID
	u.Username = user.Username
}

// FromDomain is the function that builds a http model based on slice based on the model slice's data from domain.
func (us *Users) FromDomain(users domainentity.Users) {
	u := User{}

	for _, user := range users {
		u.FromDomain(user)
		*us = append(*us, u)
	}
}
