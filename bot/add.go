package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mvdan.cc/xurls"
)

type Add struct{}

// OnMessage returns N last news articles
func (a Add) OnMessage(msg tgbotapi.Message) Response {
	if !contains(a.ReactOn(), msg.Command()) {
		return Response{}
	}

	args := msg.CommandArguments()
	url := xurls.Relaxed.FindString(args)
	title := strings.ReplaceAll(args, url+" ", "")

	return Response{
		Text: fmt.Sprintf("URL: %s\nTitle: %s", url, title),
		Send: true,
	}
}

// ReactOn keys
func (a Add) ReactOn() []string {
	return []string{"add", "new"}
}
