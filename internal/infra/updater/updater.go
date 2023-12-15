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
	contentManager      episode.ContentManager
	delay               time.Duration
}

func NewUpdater(
	subUC *subscription.SubscriptionUsecase,
	cm episode.ContentManager,
	delay time.Duration,
) Updater {
	return Updater{subscriptionUsecase: subUC, contentManager: cm, delay: delay}
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

			// TODO: May be pass Subscription entity
			_, err := u.contentManager.CheckUpdate(ctx, s.URL(), s.LastUpdated(), s.Filter())
			if err != nil {
				log.Printf("[WARN] failed to check update for %s, %+v", s.URL(), err)
				return
			}

			// TODO: Add handling
			// for _, dl := range dls {
			// 	if err = fs.saveEpisode(dl); err != nil {
			// 		log.Printf("[ERROR] failed to add episode, %v", err)
			// 		return
			// 	}
			// }

			// s.LastUpdated = time.Now()
			// if err := fs.Store.SaveSubsctiption(&s); err != nil {
			// 	log.Printf("[WARN] failed to update time, %v", err)
			// }
		}(sub)
	}

	go wg.Wait()
	log.Println("[INFO] update completed")
}
