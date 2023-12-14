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
	go func() {
		<-ctx.Done()
		log.Printf("[INFO] stopping http server...")
		if s.server != nil {
			if err := s.server.Close(); err != nil {
				log.Printf("[ERROR] failed to close http server, %+v", err)
			}
		}
	}()

	log.Printf("[INFO] starting server at %v...\n", s.port)
	return s.server.ListenAndServe()
}
