package security

// Passwords is the model of the user current and new passwords.
type Passwords struct {
	CurrentPassword string `validate:"nonzero, password" json:"current_password"`
	NewPassword     string `validate:"nonzero, password" json:"new_password"`
}
