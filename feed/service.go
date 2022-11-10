package feed

import (
	"context"
	"log"
	"time"
	"wckd1/tg-youtube-podcasts-bot/db"
	"wckd1/tg-youtube-podcasts-bot/file_manager"
)

type FeedService struct {
	Context     context.Context
	Store       db.Store
	FileManager file_manager.FileManager
}

// Add to subscriptions
func (fs FeedService) Add(arg string) error {
	sub, err := fs.parseSubscription(arg)
	if err != nil {
		log.Printf("[ERROR] failed to parse arguments, %v", err)
		return err
	}

	if sub.IsVideo {
		// If requested single video - just load it
		if err = fs.addVideo(sub); err != nil {
			log.Printf("[ERROR] failed to create subscription, %v", err)
			return err
		}
	} else {
		// Or add to subscriptions for later updates check
		if err = fs.addSubsctiption(sub); err != nil {
			log.Printf("[ERROR] failed to create subscription, %v", err)
			return err
		}
	}

	return nil
}

// Delete from subscriptions
func (fs FeedService) Delete(arg string) error {
	sub, err := fs.parseSubscription(arg)
	if err != nil {
		log.Printf("[ERROR] failed to parse arguments, %v", err)
		return err
	}

	if err := fs.Store.DeleteSubsctiption(&sub); err != nil {
		log.Printf("[ERROR] failed to remove subscription, %v", err)
		return err
	}

	if err = fs.Store.DeleteUpdate(sub.ID); err != nil {
		log.Printf("[ERROR] failed to remove update for %s, %v", sub.ID, err)
		return err
	}

	return nil
}

// Get list of available episodes
func (fs FeedService) GetEpisodes() (eps []db.Episode, err error) {
	eps, err = fs.Store.GetEpisodes(20)
	if err != nil {
		log.Printf("[ERROR] failed to get episode, %v", err)
	}
	return
}

// Get list of pending subscriptions
func (fs FeedService) GetPendingSubsctiptions() (subs []db.Subscription, err error) {
	// Get saved updates
	upds, err := fs.Store.GetUpdates()
	if err != nil {
		log.Printf("[ERROR] failed to get updates, %v", err)
		return
	}

	now := time.Now()

	for _, upd := range upds {
		// Calculate next update time for subscription
		updt := upd.LastUpdated.Add(upd.UpdateInterval)

		if updt.Before(now) || updt.Equal(now) {
			// Get subscription if update required
			sub, err := fs.Store.GetSubsctiption(upd.SubscriptionID)
			if err != nil {
				log.Printf("[ERROR] failed to get subscription, %v", err)
				continue
			}
			subs = append(subs, sub)
		}
	}

	return
}

// Handle single video request
func (fs FeedService) addVideo(sub db.Subscription) error {
	dl, err := fs.FileManager.Get(fs.Context, sub.URL)
	if err != nil {
		return err
	}

	ep := db.Episode{
		Enclosure: db.Enclosure{
			URL:    dl.URL,
			Length: dl.Info.Length,
			Type:   "audio/mpeg",
		},
		Link:        dl.Info.Link,
		Image:       dl.Info.ImageURL,
		Title:       dl.Info.Title,
		Description: "<![CDATA[" + dl.Info.Description + "]]>",
		Author:      dl.Info.Author,
		Duration:    dl.Info.Duration,
		PubDate:     fs.parseDate(dl.Info.Date),
	}
	if err = fs.Store.CreateEpisode(&ep); err != nil {
		log.Printf("[ERROR] failed to create episode, %v", err)
		return err
	}

	return nil
}

// Handle subscription request
func (fs FeedService) addSubsctiption(sub db.Subscription) error {
	err := fs.Store.CreateSubsctiption(&sub)
	if err != nil {
		log.Printf("[ERROR] failed to create subscription, %v", err)
		return err
	}

	interval, _ := time.ParseDuration("24h")

	update := db.Update{
		SubscriptionID: sub.ID,
		UpdateInterval: interval,
		LastUpdated:    time.Now(),
	}

	if err = fs.Store.ChangeUpdate(&update); err != nil {
		log.Printf("[ERROR] failed to create update, %v", err)
		return err
	}

	return nil
}
