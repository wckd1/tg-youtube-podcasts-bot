package bot

import (
	"log"
	db "wckd1/tg-youtube-podcasts-bot/db"
)

type Remove struct {
	Store   db.Store
}

// OnMessage return deleted subscription status
func (a Remove) OnMessage(msg Message) Response {
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
	err = a.Store.DeleteSubsctiption(&sub)
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
