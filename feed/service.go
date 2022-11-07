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

	switch sub.Type {
	// If requested single video - just load it
	case db.Video:
		if err = fs.addVideo(sub); err != nil {
			log.Printf("[ERROR] failed to create subscription, %v", err)
			return err
		}
	// Or add to subscriptions for later updates check
	default:
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

	if err = fs.Store.DeleteSubsctiption(&sub); err != nil {
		log.Printf("[ERROR] failed to remove subscription, %v", err)
		return err
	}

	return nil
}

// Handle single video request
func (fs FeedService) addVideo(sub db.Subscription) error {
	dl, err := fs.FileManager.Get(fs.Context, sub.YouTubeID)
	if err != nil {
		return err
	}

	ep := db.Episode{
		URL:         dl.URL,
		CoverURL:    dl.CoverURL,
		Title:       dl.Title,
		Description: dl.Description,
	}
	if err = fs.Store.CreateEpisode(&ep); err != nil {
		log.Printf("[ERROR] failed to create episode, %v", err)
		return err
	}

	return nil
}

// Handle subscription request
func (fs FeedService) addSubsctiption(sub db.Subscription) error {
	subID, err := fs.Store.CreateSubsctiption(&sub)
	if err != nil {
		log.Printf("[ERROR] failed to create subscription, %v", err)
		return err
	}

	update := db.Update{
		SubscriptionID: subID,
		UpdateInterval: "1d",
		LastUpdated:    time.Now(),
	}

	if err = fs.Store.ChangeUpdate(&update); err != nil {
		log.Printf("[ERROR] failed to create update, %v", err)
		return err
	}

	return nil
}
