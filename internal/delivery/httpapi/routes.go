package httpapi

import (
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/httpapi/handler"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"
)

// RSS
func (s *HTTPServer) RegisterRSSHandler(rssUsecase *usecase.RSSUseCase) {
	rssHandler := handler.NewRSSHandler(rssUsecase)

	s.router.Get("/rss/{playlist_id}", rssHandler.GetRSS)
}
