package entity

import (
	uuid "github.com/satori/go.uuid"
)

// Auth is the representation of the auth's domain entity.
type Auth struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

// IsEmpty is the function that checks if auth's domain entity is empty.
func (a Auth) IsEmpty() bool {
	return a == Auth{}
}
