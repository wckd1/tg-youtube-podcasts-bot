package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	db "wckd1/tg-youtube-podcasts-bot/db/store"

	"mvdan.cc/xurls"
)

type Remove struct {
	Context context.Context
	Store   db.Store
}

// OnMessage return deleted subscription status
func (a Remove) OnMessage(msg Message) Response {
	if !contains(a.ReactOn(), msg.Command) {
		return Response{}
	}

	params, err := parseRemoveParams(msg.Arguments)
	if err != nil {
		log.Printf("[ERROR] failed to parse arguments, %v", err)
		return Response{
			Text: err.Error(),
			Send: true,
		}
	}
	if len(params.Title) > 0 {
		err = a.Store.DeleteTitledSubsctiption(a.Context, params)
	} else {
		err = a.Store.DeleteSubsctiption(a.Context, params)
	}
	if err != nil {
		log.Printf("[ERROR] failed to remove subscription, %v", err)
		return Response{
			Text: "Failed to remove subscription",
			Send: true,
		}
	}

	return Response{
		Text: "Unsubscribed",
		Send: true,
	}
}

// ReactOn keys
func (a Remove) ReactOn() []string {
	return []string{"remove", "rm", "delete", "unsub"}
}

// Parse params to detect subscription details
func parseRemoveParams(arguments string) (params db.DeleteSubscriptionParams, err error) {
	params = db.DeleteSubscriptionParams{}

	// Check if arguments contains link
	furl := xurls.Relaxed.FindString(arguments)
	if len(furl) == 0 {
		err = fmt.Errorf("no source url provided")
		return
	}
	params.SourcePath = furl

	// Parse title
	title := strings.ReplaceAll(arguments, furl, "")
	params.Title = strings.TrimSpace(title)

	return
}
