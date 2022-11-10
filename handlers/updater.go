package handlers

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	"wckd1/tg-youtube-podcasts-bot/feed"
)

const (
	ytdlpCmd = "yt-dlp --dump-json %s"
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

	sl, err := u.FeedService.GetPendingSubsctiptions()
	if err != nil {
		log.Printf("[WARN] updates check skipped, %w", err)
	}

	testSub := sl[0]

	cmdStr := fmt.Sprintf(ytdlpCmd, testSub.URL)
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	log.Printf("[DEBUG] executing command: %s", cmd.String())
	if err = cmd.Run(); err != nil {
		log.Printf("[ERROR] failed to execute command: %v", err)
		return
	}
}

// All filters, separate json files
// yt-dlp --dateafter 20200520 --match-filters title~='MOUNTAIN BIKE' https://www.youtube.com/playlist\?list\=PLWx61XgoQmqdkfWC58_sYKAZvdQt9eBxQ --write-info-json --skip-download --no-write-playlist-metafiles

// For filter by title
// --match-filters title~='{title}'

// For filter by date
// --dateafter "YYYYMMDD"

// Get videos info
// -j, --dump-json
// -J, --dump-single-json
