package response_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestRespondWithJson() {
	resprec := &httptest.ResponseRecorder{}
	statusCode := 0
	var payload interface{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRespondingWithOKAndJsonBody",
			SetUp: func(t *testing.T) {
				resprec = httptest.NewRecorder()
				statusCode = http.StatusOK
				text := "everything is up and running"
				payload = messagehttputilpkg.Message{Text: text}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailInRespondingAndJsonIfItIsNotPossibleToGetJsonEncodingOfPayload",
			SetUp: func(t *testing.T) {
				resprec = httptest.NewRecorder()
				statusCode = http.StatusInternalServerError
				payload = func() {
					customerror.New("failed")
				}
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			responsehttputilpkg.RespondWithJson(resprec, statusCode, payload)

			if !tc.WantError {
				assert.Equal(t, resprec.Result().Header.Get("Content-Type"), "application/json")
				assert.Equal(t, statusCode, resprec.Result().StatusCode)
				returnedMessage := messagehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.Equal(t, payload, returnedMessage)
			} else {
				assert.Equal(t, resprec.Result().Header.Get("Content-Type"), "application/json")
				assert.Equal(t, statusCode, resprec.Result().StatusCode)
			}
		})
	}
}
