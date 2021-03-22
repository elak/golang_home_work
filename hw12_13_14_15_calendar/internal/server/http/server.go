package internalhttp

import (
	"context"
	"errors"
	"net/http"
)

type Server struct {
	Port     string
	Host     string
	Instance *http.Server
}

type Application interface{}

func NewServer(app Application, host string, port string) *Server {
	var s Server

	s.Host = host
	s.Port = port

	http.Handle("/", loggingMiddleware(http.HandlerFunc(s.HandleSayHello)))

	s.Instance = &http.Server{Addr: s.Host + ":" + s.Port}

	return &s
}

func (s *Server) HandleSayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello."))
}

func (s *Server) Start(ctx context.Context) error {
	err := s.Instance.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Instance.Shutdown(ctx)
}
