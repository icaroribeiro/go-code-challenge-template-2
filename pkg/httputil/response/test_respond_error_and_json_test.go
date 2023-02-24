package response_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/customerror"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/response"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestRespondErrorAndJson() {
	res := &httptest.ResponseRecorder{}
	statusCode := 0
	var err error
	payload := responsehttputilpkg.Error{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRespondingWithInternalServerErrorAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusInternalServerError
				text := "failed"
				err = customerror.New(text)
				payload = responsehttputilpkg.Error{Text: text}
			},
		},
		{
			Context: "ItShouldSucceedInRespondingWithBadRequestAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusBadRequest
				text := "failed"
				err = customerror.BadRequest.New(text)
				payload = responsehttputilpkg.Error{Text: text}
			},
		},
		{
			Context: "ItShouldSucceedInRespondingWithUnauthorizedAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusUnauthorized
				text := "failed"
				err = customerror.Unauthorized.New(text)
				payload = responsehttputilpkg.Error{Text: text}
			},
		},
		{
			Context: "ItShouldSucceedInRespondingWithNotFoundAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusNotFound
				text := "failed"
				err = customerror.NotFound.New(text)
				payload = responsehttputilpkg.Error{Text: text}
			},
		},
		{
			Context: "ItShouldSucceedInRespondingWithConflictAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusConflict
				text := "failed"
				err = customerror.Conflict.New(text)
				payload = responsehttputilpkg.Error{Text: text}
			},
		},
		{
			Context: "ItShouldSucceedInRespondingWithUnprocessableEntityAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusUnprocessableEntity
				text := "failed"
				err = customerror.UnprocessableEntity.New(text)
				payload = responsehttputilpkg.Error{Text: text}
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			responsehttputilpkg.RespondErrorWithJSON(res, err)

			assert.Equal(t, res.Result().Header.Get("Content-Type"), "application/json")
			assert.Equal(t, statusCode, res.Result().StatusCode)
			errMessage := responsehttputilpkg.Error{}
			err := json.NewDecoder(res.Body).Decode(&errMessage)
			assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			assert.Equal(t, payload, errMessage)
		})
	}
}
