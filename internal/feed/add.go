package feed

import (
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/db"
	"wckd1/tg-youtube-podcasts-bot/internal/file_manager"
)

// Handle single video request
func (fs FeedService) addEpisode(sub db.Subscription) error {
	dl, err := fs.FileManager.Get(fs.Context, sub.URL)
	if err != nil {
		return err
	}

	return fs.saveEpisode(dl)
}

// Save uploaded episode to storage
func (fs FeedService) saveEpisode(dl file_manager.Download) error {
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

	return fs.Store.CreateEpisode(&ep)
}

// Handle subscription request
func (fs FeedService) addSubsctiption(sub db.Subscription) error {
	interval, _ := time.ParseDuration("24h")

	sub.UpdateInterval = interval
	sub.LastUpdated = time.Now()

	return fs.Store.SaveSubsctiption(&sub)
}
