package feed

import (
	"log"
	"time"
	"wckd1/tg-youtube-podcasts-bot/db"
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
	for _, sub := range pendings {
		// TODO: Run in gorutines with channel to handle finish
		dls, err := fs.FileManager.CheckUpdate(fs.Context, sub.URL, sub.LastUpdated, sub.Filter)
		if err != nil {
			log.Printf("[WARN] failed to check update, %v", err)
			continue
		}

		for _, dl := range dls {
			if err = fs.saveEpisode(dl); err != nil {
				log.Printf("[ERROR] failed to add episode, %v", err)
				continue
			}
		}

		sub.LastUpdated = time.Now()
		if err := fs.Store.SaveSubsctiption(&sub); err != nil {
			log.Printf("[WARN] failed to update time, %v", err)
		}
	}
}
