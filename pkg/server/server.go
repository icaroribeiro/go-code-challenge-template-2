package server

import (
	"context"
	"net/http"
)

type Server struct {
	Instance *http.Server
}

// New is the factory function that encapsulates the implementation related to server.
func New(tcpAddress string, handler http.Handler) *Server {
	return &Server{
		Instance: &http.Server{
			Addr:    tcpAddress,
			Handler: handler,
		},
	}
}

// Start is the function that starts the server.
func (s *Server) Start() error {
	return s.Instance.ListenAndServe()
}

// Stop is the function that stops the server.
func (s *Server) Stop(ctx context.Context) error {
	return s.Instance.Shutdown(ctx)
}
