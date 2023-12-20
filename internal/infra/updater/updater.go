package updater

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/service"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"
)

var (
	ErrContextClosed = errors.New("context closed")
)

// Updater is a task runner that check for updates with given delay
type Updater struct {
	updateUsecase  *usecase.UpdateUsecase
	contentManager service.ContentManager
	ticker         time.Ticker
	delay          time.Duration
}

func NewUpdater(
	updateUsecase *usecase.UpdateUsecase,
	contentManager service.ContentManager,
	delay time.Duration,
) Updater {
	ticker := time.NewTicker(delay)
	return Updater{updateUsecase, contentManager, *ticker, delay}
}

func (u Updater) Start(ctx context.Context) {
	log.Printf("[INFO] check for updates on startup")
	u.checkForUpdates(ctx)

	for range u.ticker.C {
		log.Printf("[INFO] checking for updates...")
		u.checkForUpdates(ctx)
	}
}

func (u Updater) Shutdown() {
	u.ticker.Stop()
	log.Println("[INFO] updater stopped")
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
		go func(s entity.Subscription) {
			defer wg.Done()

			eps, err := u.contentManager.CheckUpdate(ctx, s)
			if err != nil {
				log.Printf("[WARN] failed to check update for %s, %+v", s.URL(), err)
				return
			}

			wg.Add(len(eps))
			for _, ep := range eps {
				go func(e entity.Episode) {
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
