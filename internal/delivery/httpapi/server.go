package httpapi

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTPServer struct {
	port   string
	server *http.Server
	router chi.Router
}

func NewServer(port string) *HTTPServer {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	return &HTTPServer{
		port:   port,
		server: srv,
		router: r,
	}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	log.Printf("[INFO] starting server at %v...\n", s.port)
	return s.server.ListenAndServe()
}

func (s HTTPServer) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	log.Printf("[INFO] http server stopped")
	return nil
}
