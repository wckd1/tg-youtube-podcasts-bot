package feed

import (
	"log"
	"time"
)

func (fs FeedService) CheckForUpdates() {
	subs, err := fs.Store.GetSubscriptions()
	if err != nil {
		log.Printf("[WARN] updates check skipped, %v", err)
	}

	now := time.Now()

	for _, sub := range subs {
		// Calculate next update time for subscription
		updt := sub.LastUpdated.Add(sub.UpdateInterval)

		if updt.Before(now) || updt.Equal(now) {
			// Update subscription if required
			// TODO: Run in gorutines with channel to handle finish
			lu := sub.LastUpdated.Format("20060102") // May be -1
			
			dls, err := fs.FileManager.CheckUpdate(fs.Context, sub.URL, lu)
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
}
