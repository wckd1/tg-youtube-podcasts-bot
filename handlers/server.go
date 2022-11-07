package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"wckd1/tg-youtube-podcasts-bot/rss"
)

// Server provides HTTP API
type Server struct {
	httpServer *http.Server
}

// Run starts http server for API
func (s *Server) Run(ctx context.Context, port int) {
	go func() {
		<-ctx.Done()
		log.Printf("[INFO] stopping http server")
		if s.httpServer != nil {
			if clsErr := s.httpServer.Close(); clsErr != nil {
				log.Printf("[ERROR] failed to close proxy http server, %v", clsErr)
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/rss", getRSS)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	
	log.Printf("[INFO] starting server at %d", port)
	s.httpServer.ListenAndServe()
}

func getRSS(w http.ResponseWriter, r *http.Request) {
	rss, err := rss.RSSFeed()
	if err != nil {
		log.Printf("[ERROR] failed generate rss feed, %v", err)

	}

	w.Header().Set("Content-Type", "application/xml; charset=UTF-8")

	io.WriteString(w, rss)
}
