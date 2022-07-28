package httpServer

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

type Server struct {
	notify chan error
	server *http.Server
}

func New(handler http.Handler, port string, readTimeout, writeTimeout time.Duration) *Server {
	s := &Server{
		notify: make(chan error, 1),
		server: &http.Server{
			Handler:      handler,
			Addr:         net.JoinHostPort("", port),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
	return s
}

func (s *Server) Run() {
	go func() {
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err
			close(s.notify)
		}
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

