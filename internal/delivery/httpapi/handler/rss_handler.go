package handler

import (
	"errors"
	"log"
	"net/http"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"

	"github.com/go-chi/chi/v5"
)

var ErrNoPlaylistID = errors.New("no playlist id is specified")

type rssHandler struct {
	rssUsecase *usecase.RSSUseCase
}

func NewRSSHandler(rssUsecase *usecase.RSSUseCase) rssHandler {
	return rssHandler{rssUsecase}
}

func (h rssHandler) GetRSS(w http.ResponseWriter, r *http.Request) {
	plID := chi.URLParam(r, "playlist_id")

	rss, err := h.rssUsecase.BuildRSS(plID)
	if err != nil {
		log.Printf("[ERROR] failed to generate rss %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rss)) // io.WriteString(w, rss)
}
