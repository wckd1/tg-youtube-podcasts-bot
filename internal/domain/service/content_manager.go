package service

import (
	"context"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
)

type ContentManager interface {
	Get(ctx context.Context, url string) (ep episode.Episode, err error)
	CheckUpdate(ctx context.Context, sub subscription.Subscription) (eps []episode.Episode, err error)
}
