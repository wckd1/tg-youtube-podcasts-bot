package bot

import (
	"wckd1/tg-youtube-podcasts-bot/feed"
)

type Remove struct {
	FeedService feed.FeedService
}

// OnMessage return deleted subscription status
func (a Remove) OnMessage(msg Message) Response {
	if !contains(a.ReactOn(), msg.Command) {
		return Response{}
	}

	if err := a.FeedService.Delete(msg.Arguments); err != nil {
		return Response{
			Text: "Failed to add subscription. See logs for more info",
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
