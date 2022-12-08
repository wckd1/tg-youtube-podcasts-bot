package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"wckd1/tg-youtube-podcasts-bot/internal/feed"
	"wckd1/tg-youtube-podcasts-bot/internal/rss"
)

// Server provides HTTP API
type Server struct {
	FeedService feed.FeedService
	httpServer  *http.Server
}

// Run starts http server for API
func (s *Server) Run(ctx context.Context, rssKey string, port int) {
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
	mux.HandleFunc("/rss/"+rssKey, s.getRSS)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("[INFO] starting server at %d", port)
	s.httpServer.ListenAndServe()
}

// Get formatted rss string
func (s *Server) getRSS(w http.ResponseWriter, r *http.Request) {
	el, err := s.FeedService.GetEpisodes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal error"))
		return
	}

	rss, err := rss.RSSFeed(el)
	if err != nil {
		log.Printf("[ERROR] failed generate rss feed, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal error"))
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=UTF-8")

	io.WriteString(w, rss)
}
