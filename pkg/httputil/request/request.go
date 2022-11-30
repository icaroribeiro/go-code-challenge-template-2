package request

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// RequestData is the model of a request data.
type RequestData struct {
	Method  string
	Target  string
	Body    interface{}
	Headers map[string][]string
}

type ContextKeyType string

type contextKey struct {
	name string
}

// PrepareRequestBody is the fucntion that formats the request body before executing a request.
func PrepareRequestBody(body interface{}) io.Reader {
	var reqBody io.Reader

	bodyString, ok := body.(string)
	if ok {
		if bodyString == "" {
			return nil
		}

		formattedBody := bodyString
		escapeSequencies := []string{"\t", "\n"}

		for _, value := range escapeSequencies {
			formattedBody = strings.ReplaceAll(formattedBody, value, "")
		}

		reqBody = strings.NewReader(formattedBody)
	}

	bodyBufferOfBytes, ok := body.(*bytes.Buffer)
	if ok {
		reqBody = bodyBufferOfBytes
	}

	return reqBody
}

// SetRequestHeaders is the function that configures the header entries before executing a request.
func SetRequestHeaders(r *http.Request, headers map[string][]string) {
	for key, values := range headers {
		for _, value := range values {
			r.Header.Set(key, value)
		}
	}
}
