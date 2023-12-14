package httpapi

import (
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/httpapi/handler"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/rss"
)

// RSS
func (s *HTTPServer) RegisterRSSHandler(rssUsecase *rss.RSSUseCase, secret string) {
	rssHandler := handler.NewRSSHandler(rssUsecase, secret)

	s.router.Get("/rss/{secret}", rssHandler.GetRSS)
}
