package episode

import (
	"context"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
)

type ContentManager interface {
	Get(ctx context.Context, url string) (ep Episode, err error)
	CheckUpdate(ctx context.Context, sub subscription.Subscription) (eps []Episode, err error)
}
