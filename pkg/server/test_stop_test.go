package server_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/mux"
	serverpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/server"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestStop() {
	tcpAddress := fmt.Sprintf(":%s", fake.Digit())

	handler := mux.NewRouter()
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler.Name("/").Methods(http.MethodGet).Path("Home").HandlerFunc(handlerFunc)

	var err error
	serverRunning := make(chan struct{})
	serverDone := make(chan struct{})

	ts.Cases = Cases{
		{
			Context:   "ItShouldSucceedInStoppingTheServer",
			SetUp:     func(t *testing.T) {},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			server := serverpkg.New(tcpAddress, handler)

			go func() {
				close(serverRunning)
				server.Start()
				defer close(serverDone)
			}()

			<-serverRunning

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			err = server.Stop(ctx)

			<-serverDone

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			}
		})
	}
}
