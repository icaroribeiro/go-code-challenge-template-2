package security

// ISecurity interface is a collection of function signatures that represents the security's contract.
type ISecurity interface {
	HashPassword(password string, cost int) ([]byte, error)
	VerifyPasswords(hashedPassword, password string) error
}
