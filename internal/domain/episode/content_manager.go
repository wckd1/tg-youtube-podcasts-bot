package episode

import (
	"context"
	"time"
)

type ContentManager interface {
	Get(ctx context.Context, id, url string) (ep Episode, err error)
	CheckUpdate(ctx context.Context, url string, date time.Time, filter string) (eps []Episode, err error)
}
