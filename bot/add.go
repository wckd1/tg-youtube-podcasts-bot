package bot

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	db "wckd1/tg-youtube-podcasts-bot/db"
	"wckd1/tg-youtube-podcasts-bot/loader"

	"mvdan.cc/xurls/v2"
)

type Add struct {
	Store  db.Store
	Loader loader.Interface
}

// OnMessage return new subscription status
func (a Add) OnMessage(msg Message) Response {
	if !contains(a.ReactOn(), msg.Command) {
		return Response{}
	}

	sub, err := parseSubscription(msg.Arguments)
	if err != nil {
		log.Printf("[ERROR] failed to parse arguments, %v", err)
		return Response{
			Text: err.Error(),
			Send: true,
		}
	}
	// If requested single video - just load it
	if sub.SourceType == db.Video {
		go a.Loader.Download(sub.SourcePath)
		return Response{
			Text: "Audio will be available shortly",
			Send: true,
		}
	}

	err = a.Store.CreateSubsctiption(&sub)
	if err != nil {
		log.Printf("[ERROR] failed to create subscription, %v", err)

		errMsg := "Failed to create subscription"
		// if sqliteErr, ok := err.(sqlite3.Error); ok {
		// 	if sqliteErr.Code == sqlite3.ErrConstraint {
		// 		errMsg = "This subscription is alreay added"
		// 	}
		// }
		return Response{
			Text: errMsg,
			Send: true,
		}
	}

	return Response{
		Text: "Subscribed",
		Send: true,
	}
}

// ReactOn keys
func (a Add) ReactOn() []string {
	return []string{"add", "new", "sub"}
}

// Parse params to detect subscription details
func parseSubscription(arguments string) (sub db.Subscription, err error) {
	sub = db.Subscription{}

	// Check if arguments contains link
	furl := xurls.Relaxed().FindString(arguments)
	if len(furl) == 0 {
		err = fmt.Errorf("no source url provided")
		return
	}

	// Check if YouTube link
	purl, err := url.Parse(furl)
	if strings.ReplaceAll(purl.Host, "www.", "") != "youtube.com" {
		err = fmt.Errorf("only youtube links are supported")
		return
	}

	// Check link type
	path := strings.Split(purl.Path, "/")[1]
	switch path {
	// Check if passed video link
	case "watch":
		// Check if video is in playlist
		if _, ok := purl.Query()["list"]; ok {
			err = fmt.Errorf("playlist subscription is not implemented yet")
			// sub.ID = listID[0]
			// sub.SourceType = db.Playlist
			// sub.SourcePath = listID[0]
		} else if videoID, ok := purl.Query()["v"]; ok {
			sub.ID = videoID[0]
			sub.SourceType = db.Video
			sub.SourcePath = videoID[0]
		} else {
			err = fmt.Errorf("unrecognized link type")
		}
	// Check if passed channel link
	case "c", "channel":
		err = fmt.Errorf("channel subscription is not implemented yet")
		// sub.SourceType = db.Channel

		// // Parse title
		// title := strings.ReplaceAll(arguments, furl, "")
		// sub.Title = strings.TrimSpace(title)

		// channelID := strings.Split(purl.Path, "/")[2]
		// sub.ID = strings.Join([]string{channelID, sub.Title}, "_")
		// sub.SourcePath = channelID
	// Other links are unsupported
	default:
		err = fmt.Errorf("unrecognized link type")
	}

	return
}
