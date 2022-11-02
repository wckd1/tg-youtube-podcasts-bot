package bot

import (
	"fmt"
	"strings"

	"mvdan.cc/xurls"
)

type Add struct{}

// OnMessage returns N last news articles
func (a Add) OnMessage(msg Message) Response {
	if !contains(a.ReactOn(), msg.Command) {
		return Response{}
	}

	url := xurls.Relaxed.FindString(msg.Arguments)
	title := strings.ReplaceAll(msg.Arguments, url+" ", "")

	return Response{
		Text: fmt.Sprintf("URL: %s\nTitle: %s", url, title),
		Send: true,
	}
}

// ReactOn keys
func (a Add) ReactOn() []string {
	return []string{"add", "new"}
}
