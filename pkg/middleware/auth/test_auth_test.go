package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/customerror"
	adapterhttputilpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/httputil/adapter"
	requesthttputilpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/httputil/response"
	routehttputilpkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template-2/pkg/middleware/auth"
	mockauthpkg "github.com/icaroribeiro/go-code-challenge-template-2/tests/mocks/pkg/mockauth"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestAuth() {
	bearerToken := []string{"", ""}

	tokenString := ""

	authHeaderString := ""

	token := &jwt.Token{}

	headers := make(map[string][]string)

	statusCode := 0
	payload := responsehttputilpkg.Message{}

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithAuthenticationMiddleware",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()

				bearerToken = []string{"Bearer", tokenString}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				authHeaderString = value
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}

				statusCode = http.StatusOK

				payload = responsehttputilpkg.Message{Text: "ok"}

				returnArgs = ReturnArgs{
					{tokenString, nil},
					{token, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheAuthorizationHeaderIsNotSetInTheRequestHeader",
			SetUp: func(t *testing.T) {
				tokenString = ""

				bearerToken = []string{"", tokenString}

				authHeaderString = ""

				headers = map[string][]string{}

				token = &jwt.Token{}

				statusCode = http.StatusOK

				returnArgs = ReturnArgs{
					{tokenString, nil},
					{token, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheTokenIsNotDecoded",
			SetUp: func(t *testing.T) {
				tokenString = fake.Word()

				bearerToken = []string{"Bearer", tokenString}

				key := "Authorization"
				value := strings.Join(bearerToken[:], " ")
				authHeaderString = value
				headers = map[string][]string{
					key: {value},
				}

				token = &jwt.Token{}

				statusCode = http.StatusUnauthorized

				returnArgs = ReturnArgs{
					{tokenString, nil},
					{token, customerror.New("failed")},
				}
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authN := new(mockauthpkg.Auth)
			authN.On("ExtractTokenString", authHeaderString).Return(returnArgs[0]...)
			authN.On("DecodeToken", bearerToken[1]).Return(returnArgs[1]...)

			authMiddleware := authmiddlewarepkg.Auth(authN)

			handlerFunc := func(w http.ResponseWriter, r *http.Request) {
				responsehttputilpkg.RespondWithJSON(w, http.StatusOK, responsehttputilpkg.Message{Text: "ok"})
			}

			returnedHandlerFunc := adapterhttputilpkg.AdaptFunc(handlerFunc).With(authMiddleware)

			route := routehttputilpkg.Route{
				Name:        "Testing",
				Method:      http.MethodGet,
				Path:        "/testing",
				HandlerFunc: returnedHandlerFunc,
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			requesthttputilpkg.SetRequestHeaders(req, headers)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Result().Header.Get("Content-Type"), "application/json")
				assert.Equal(t, statusCode, resprec.Result().StatusCode)
				returnedMessage := responsehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.Equal(t, payload, returnedMessage)
			} else {
				assert.Equal(t, statusCode, resprec.Result().StatusCode)
			}
		})
	}
}
