package bot

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	db "wckd1/tg-youtube-podcasts-bot/db/store"

	"github.com/mattn/go-sqlite3"
	"mvdan.cc/xurls"
)

type Add struct {
	Context context.Context
	Store   db.Store
}

// OnMessage return new subscription status
func (a Add) OnMessage(msg Message) Response {
	if !contains(a.ReactOn(), msg.Command) {
		return Response{}
	}

	params, err := parseParams(msg.Arguments)
	if err != nil {
		log.Printf("[ERROR] failed to parse arguments, %v", err)
		return Response{
			Text: err.Error(),
			Send: true,
		}
	}

	err = a.Store.CreateSubsctiption(a.Context, params)
	if err != nil {
		log.Printf("[ERROR] failed to create subscription, %v", err)

		errMsg := "Failed to create subscription"
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				errMsg = "This subscription is alreay added"
			}
		}
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
	return []string{"add", "new"}
}

// Parse params to detect subscription details
func parseParams(arguments string) (params db.CreateSubscriptionParams, err error) {
	params = db.CreateSubscriptionParams{}

	// Check if arguments contains link
	furl := xurls.Relaxed.FindString(arguments)
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
	params.SourcePath = furl

	// Check link type
	path := strings.Split(purl.Path, "/")[1]
	switch path {
	// Check if passed video link
	case "watch":
		// Check if video is in playlist
		if _, ok := purl.Query()["list"]; ok {
			params.SourceType = db.Playlist
		} else {
			params.SourceType = db.Video
		}
	// Check if passed channel link
	case "c", "channel":
		params.SourceType = db.Channel

		// Parse title
		title := strings.ReplaceAll(arguments, furl, "")
		params.Title = strings.TrimSpace(title)
	// Other links are unsupported
	default:
		err = fmt.Errorf("unrecognized link type")
	}

	return
}
