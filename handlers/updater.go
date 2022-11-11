package handlers

import (
	"context"
	"log"
	"time"
	"wckd1/tg-youtube-podcasts-bot/feed"
)

// Updater is a task runner that check for updates with given delay
type Updater struct {
	Delay       time.Duration
	FeedService feed.FeedService
}

func (u Updater) Start(ctx context.Context) {
	log.Printf("[INFO] starting updater with %v interval", u.Delay)

	ticker := time.NewTicker(u.Delay)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			log.Printf("[INFO] Check for updates")
			u.FeedService.CheckForUpdates()

		case <-ctx.Done():
			log.Printf("[INFO] update checker stopped, %v", ctx.Err())
			return
		}
	}
}
