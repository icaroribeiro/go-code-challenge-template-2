package request

import "net/http"

// RequestData is the model of a request data.
type RequestData struct {
	Method string
	Target string
	Body   interface{}
}

// SetRequestHeaders is the function that configures the header entries before executing a request.
func SetRequestHeaders(r *http.Request, headers map[string][]string) {
	for key, values := range headers {
		for _, value := range values {
			r.Header.Set(key, value)
		}
	}
}
