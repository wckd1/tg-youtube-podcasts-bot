package feed

import (
	"log"
	"sync"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/db"
)

func (fs FeedService) CheckForUpdates() {
	subs, err := fs.Store.GetSubscriptions()
	if err != nil {
		log.Printf("[WARN] updates check skipped, %v", err)
		return
	}

	now := time.Now()

	var pendings []db.Subscription
	for _, sub := range subs {
		// Calculate next update time for subscription
		updt := sub.LastUpdated.Add(sub.UpdateInterval)

		if updt.Before(now) || updt.Equal(now) {
			pendings = append(pendings, sub)
		}
	}

	if len(pendings) == 0 {
		log.Printf("[Info] no updates are required")
	}

	// Update subscription if required
	wg := sync.WaitGroup{}
	wg.Add(len(pendings))

	for _, sub := range pendings {
		go func(s db.Subscription) {
			defer wg.Done()

			dls, err := fs.FileManager.CheckUpdate(fs.Context, s.URL, s.LastUpdated, s.Filter)
			if err != nil {
				log.Printf("[WARN] failed to check update, %v", err)
				return
			}

			for _, dl := range dls {
				if err = fs.saveEpisode(dl); err != nil {
					log.Printf("[ERROR] failed to add episode, %v", err)
					return
				}
			}

			s.LastUpdated = time.Now()
			if err := fs.Store.SaveSubsctiption(&s); err != nil {
				log.Printf("[WARN] failed to update time, %v", err)
			}
		}(sub)
	}

	go wg.Wait()

	log.Println("[INFO] update completed")
}
