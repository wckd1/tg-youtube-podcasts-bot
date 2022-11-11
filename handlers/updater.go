package handlers

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	"wckd1/tg-youtube-podcasts-bot/feed"
	"wckd1/tg-youtube-podcasts-bot/db"
)

const (
	ytdlpCmd = "yt-dlp %s --skip-download --write-info-json --no-write-playlist-metafiles --dateafter %s"
	titleFilter = "--match-filters title~='%s'"
	destPath = "./storage/downloads/"
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
			u.checkForUpdates(ctx)

		case <-ctx.Done():
			log.Printf("[INFO] update checker stopped, %v", ctx.Err())
			return
		}
	}
}

func (u Updater) checkForUpdates(ctx context.Context) {
	log.Printf("[INFO] Check for updates")

	subs, err := u.FeedService.GetSubscriptions()
	if err != nil {
		log.Printf("[WARN] updates check skipped, %w", err)
	}

	now := time.Now()

	for _, sub := range subs {
		// Calculate next update time for subscription
		updt := sub.LastUpdated.Add(sub.UpdateInterval)

		if updt.Before(now) || updt.Equal(now) {
			// Update subscription if required
			// TODO: Run in gorutines with channel to handle finish
			u.updateSubscription(ctx, sub)
		}
	}
}

func (u Updater) updateSubscription(ctx context.Context, sub db.Subscription) {
	date := sub.LastUpdated.Format("20060102") // May be -1
	cmdStr := fmt.Sprintf(ytdlpCmd, sub.URL, date)
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Dir = destPath

	log.Printf("[DEBUG] executing command: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] failed to execute command: %v", err)
		return
	}

	// TODO: Pass to FileManager
}

// All filters, separate json files
// yt-dlp --dateafter 20200520 --match-filters title~='MOUNTAIN BIKE' https://www.youtube.com/playlist\?list\=PLWx61XgoQmqdkfWC58_sYKAZvdQt9eBxQ --write-info-json --skip-download --no-write-playlist-metafiles

// For filter by title
// --match-filters title~='{title}'

// For filter by date
// --dateafter "YYYYMMDD"
