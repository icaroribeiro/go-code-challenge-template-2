package validator

// IValidator transport is the validator's contract.
type IValidator interface {
	Validate(i interface{}) error
	ValidateWithTags(i interface{}, tags string) error
}
