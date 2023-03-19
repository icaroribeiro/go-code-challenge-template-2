package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
)

// RespondWithJSON is the function that generates a JSON response along with the suitable header and given status code.
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to get the JSON encoding of %+v: %s", payload, err.Error())))
		return
	}

	w.Write(response)
}

// RespondErrorWithJSON is the function that generates a JSON error response.
func RespondErrorWithJSON(w http.ResponseWriter, err error) {
	statusCode := 0

	errorType := customerror.GetType(err)

	switch errorType {
	case customerror.BadRequest:
		statusCode = http.StatusBadRequest
	case customerror.Unauthorized:
		statusCode = http.StatusUnauthorized
	case customerror.NotFound:
		statusCode = http.StatusNotFound
	case customerror.Conflict:
		statusCode = http.StatusConflict
	case customerror.UnprocessableEntity:
		statusCode = http.StatusUnprocessableEntity
	default:
		statusCode = http.StatusInternalServerError
	}

	RespondWithJSON(w, statusCode, Error{Text: err.Error()})
}
