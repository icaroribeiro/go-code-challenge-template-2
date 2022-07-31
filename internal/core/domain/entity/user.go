package entity

import (
	uuid "github.com/satori/go.uuid"
)

// User is the representation of the user's domain entity.
type User struct {
	ID       uuid.UUID
	Username string
}

// Users is a slice of User.
type Users []User
