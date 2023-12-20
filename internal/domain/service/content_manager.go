package service

import (
	"context"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
)

type ContentManager interface {
	Get(ctx context.Context, url string) (ep entity.Episode, err error)
	CheckUpdate(ctx context.Context, sub entity.Subscription) (eps []entity.Episode, err error)
}
