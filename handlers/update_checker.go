package handlers

import (
	"context"
	"time"
	"wckd1/tg-youtube-podcasts-bot/loader"
)

// UpdateChecker is a task runner that check for updates with delay
type UpdateChecker struct {
	Delay     time.Duration
	Submitter Submitter
	Loader    loader.Interface
}

// Submitter defines interface to submit to the chat
type Submitter interface {
	SubmitText(ctx context.Context, text string)
}

func (uc UpdateChecker) Start(ctx context.Context) error {
	ticker := time.NewTicker(uc.Delay)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			uc.Submitter.SubmitText(ctx, "Checked for update")

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
