package updater

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
)

var (
	ErrContextClosed = errors.New("context closed")
)

// Updater is a task runner that check for updates with given delay
type Updater struct {
	subscriptionUsecase *subscription.SubscriptionUsecase
	episodeUsecase      *episode.EpisodeUsecase
	contentManager      episode.ContentManager
	delay               time.Duration
}

func NewUpdater(
	subUC *subscription.SubscriptionUsecase,
	epUC *episode.EpisodeUsecase,
	cm episode.ContentManager,
	delay time.Duration,
) Updater {
	return Updater{subscriptionUsecase: subUC, episodeUsecase: epUC, contentManager: cm, delay: delay}
}

func (u Updater) Start(ctx context.Context) {
	log.Printf("[INFO] check for updates on startup")
	u.checkForUpdates(ctx)

	log.Printf("[INFO] starting updater with %v interval", u.delay)
	ticker := time.NewTicker(u.delay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("[INFO] checking for updates...")
			u.checkForUpdates(ctx)

		case <-ctx.Done():
			log.Printf("[INFO] context closed, %v", ctx.Err())
			return
		}
	}
}

func (u Updater) checkForUpdates(ctx context.Context) {
	subs, err := u.subscriptionUsecase.GetPendingSubscriptions()
	if err != nil {
		log.Printf("[WARN] updates check skipped, %+v", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(subs))

	for _, sub := range subs {
		go func(s subscription.Subscription) {
			defer wg.Done()

			eps, err := u.contentManager.CheckUpdate(ctx, s)
			if err != nil {
				log.Printf("[WARN] failed to check update for %s, %+v", s.URL(), err)
				return
			}

			wg.Add(len(eps))
			for _, ep := range eps {
				go func(e episode.Episode) {
					defer wg.Done()

					// TODO: Update bounded playlists
					if err = u.episodeUsecase.SaveEpisode(&e); err != nil {
						log.Printf("[ERROR] failed to add episode, %v", err)
						return
					}
				}(ep)
			}

			s.SetLastUpdated(time.Now())
			if err := u.subscriptionUsecase.SaveSubsctiption(&s); err != nil {
				log.Printf("[WARN] failed to update time, %v", err)
			}

		}(sub)
	}

	wg.Wait()
	log.Println("[INFO] update completed")
}
