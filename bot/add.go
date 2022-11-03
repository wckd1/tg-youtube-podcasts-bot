package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	db "wckd1/tg-youtube-podcasts-bot/db/store"

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

	url := xurls.Relaxed.FindString(msg.Arguments)
	title := strings.ReplaceAll(msg.Arguments, url+" ", "")

	params := db.CreateSubscriptionParams{
		Channel: url,
		Title:   title,
	}
	err := a.Store.CreateSubsctiption(a.Context, params)
	if err != nil {
		log.Printf("[ERROR] failed to create subscription, %v", err)
		return Response{
			Text: fmt.Sprintf("Failed to create subscription for %s", title),
			Send: true,
		}
	}

	return Response{
		Text: fmt.Sprintf("Subscribed to  %s", title),
		Send: true,
	}
}

// ReactOn keys
func (a Add) ReactOn() []string {
	return []string{"add", "new"}
}
