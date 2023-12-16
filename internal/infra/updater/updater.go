package updater

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/service"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"
)

var (
	ErrContextClosed = errors.New("context closed")
)

// Updater is a task runner that check for updates with given delay
type Updater struct {
	updateUsecase  *usecase.UpdateUsecase
	contentManager service.ContentManager
	delay          time.Duration
}

func NewUpdater(
	updateUsecase *usecase.UpdateUsecase,
	contentManager service.ContentManager,
	delay time.Duration,
) Updater {
	return Updater{updateUsecase, contentManager, delay}
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
	subs, err := u.updateUsecase.GetPendingSubscriptions()
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

					if err = u.updateUsecase.SaveEpisode(&e, s.ID()); err != nil {
						log.Printf("[ERROR] failed to add episode, %v", err)
						return
					}
				}(ep)
			}

			if err := u.updateUsecase.SaveSubsctiption(&s, time.Now()); err != nil {
				log.Printf("[WARN] failed to update time, %v", err)
			}

		}(sub)
	}

	wg.Wait()
	log.Println("[INFO] update completed")
}
