package feed

import (
	"context"
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/db"
	"wckd1/tg-youtube-podcasts-bot/internal/file_manager"
)

type FeedService struct {
	Context     context.Context
	Limit       int
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
		if err = fs.addEpisode(sub); err != nil {
			log.Printf("[ERROR] failed to add episode, %v", err)
			return err
		}
	} else {
		// Or add to subscriptions for later updates check
		if err = fs.addSubsctiption(sub); err != nil {
			log.Printf("[ERROR] failed to add subscription, %v", err)
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

	return nil
}

// Get list of available episodes
func (fs FeedService) GetEpisodes() (eps []db.Episode, err error) {
	eps, err = fs.Store.GetEpisodes(fs.Limit)
	if err != nil {
		log.Printf("[ERROR] failed to get episode, %v", err)
	}
	return
}
