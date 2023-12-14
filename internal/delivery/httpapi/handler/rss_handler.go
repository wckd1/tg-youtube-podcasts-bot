package handler

import (
	"log"
	"net/http"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/rss"
)

type rssHandler struct {
	rssUsecase *rss.RSSUseCase
	secret     string
}

func NewRSSHandler(rssUsecase *rss.RSSUseCase, secret string) rssHandler {
	return rssHandler{rssUsecase, secret}
}

func (h rssHandler) GetRSS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	secret, ok := ctx.Value("secret").(string)
	if !ok || secret != h.secret {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rss, err := h.rssUsecase.BuildRSS()
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
