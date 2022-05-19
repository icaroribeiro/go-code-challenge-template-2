package security

// Credentials is the model of the user credentials.
type Credentials struct {
	Username string `validate:"nonzero, username" json:"username"`
	Password string `validate:"nonzero, password" json:"password"`
}
