package error

// Error is the error's model used for handling the JSON message for an unsuccessful operation.
type Error struct {
	Text string `json:"error"`
}
