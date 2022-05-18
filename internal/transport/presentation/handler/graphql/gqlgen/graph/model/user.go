package model

import (
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	uuid "github.com/satori/go.uuid"
)

// User is the representation of user's http model.
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// Users is a slice of User.
type Users []User

// type Users []*User

// FromDomain is the function that builds a http model based on the model's data from domain.
func (u *User) FromDomain(user domainmodel.User) {
	u.ID = user.ID
	u.Username = user.Username
}

// FromDomain is the function that builds a http model based on slice based on the model slice's data from domain.
func (us *Users) FromDomain(users domainmodel.Users) {
	u := User{}

	for _, user := range users {
		u.FromDomain(user)
		*us = append(*us, u)
	}
}
