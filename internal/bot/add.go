package bot

import (
	"wckd1/tg-youtube-podcasts-bot/internal/feed"
)

type Add struct {
	FeedService feed.FeedService
}

// OnMessage return new subscription status
func (a Add) OnMessage(msg Message) Response {
	if !contains(a.ReactOn(), msg.Command) {
		return Response{}
	}

	// Add subscription
	if err := a.FeedService.Add(msg.Arguments); err != nil {
		return Response{
			Text: "Failed to add subscription. See logs for more info",
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
