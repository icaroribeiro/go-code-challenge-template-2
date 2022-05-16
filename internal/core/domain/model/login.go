package model

import (
	uuid "github.com/satori/go.uuid"
)

// Login is the representation of the login's domain model.
type Login struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Username string `create:"nonzero, username" update:"username"`
	Password string `create:"nonzero, password" update:"password"`
}

// IsEmpty is the function that checks if login's database model is empty.
func (l Login) IsEmpty() bool {
	return l == Login{}
}
