package handlers

import (
	"context"
	"log"
	"time"
)

// UpdateChecker is a task runner that check for updates with delay
type UpdateChecker struct {
	Delay     time.Duration
	Submitter Submitter
}

// Submitter defines interface to submit (usually asynchronously) to the chat
type Submitter interface {
	Submit(ctx context.Context, text string) error
}

func (uc UpdateChecker) Start(ctx context.Context) error {
	ticker := time.NewTicker(uc.Delay)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			if err := uc.Submitter.Submit(ctx, "Checked for update"); err != nil {
				log.Panic(err)
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
