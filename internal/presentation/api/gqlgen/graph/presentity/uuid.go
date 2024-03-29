package presentity

import (
	"errors"

	"github.com/99designs/gqlgen/graphql"
	uuid "github.com/satori/go.uuid"
)

// MarshalUUID is the function that allows uuid to be marshalled by graphql.
func MarshalUUID(id uuid.UUID) graphql.Marshaler {
	return graphql.MarshalString(id.String())
}

// UnmarshalUUID is the function that allows uuid to be unmarshalled by graphql.
func UnmarshalUUID(v interface{}) (uuid.UUID, error) {
	idAsString, ok := v.(string)
	if !ok {
		return uuid.Nil, errors.New("id should be a valid UUID")
	}
	return uuid.FromString(idAsString)
}
