package bot

import (
	"context"
	"sort"
	"strings"

	"github.com/go-pkgz/syncs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command interface {
	OnMessage(msg tgbotapi.Message) Response
	ReactOn() (res []string)
}

type Response struct {
	Text string
	Send bool
}

// Commands combines many commands to one virtual
type Commands []Command

// OnMessage pass msg to all commands and collects responses (combining all of them)
func (c Commands) OnMessage(msg tgbotapi.Message) Response {
	resps := make(chan string)

	wg := syncs.NewSizedGroup(4)
	for _, command := range c {
		command := command

		wg.Go(func(ctx context.Context) {
			if resp := command.OnMessage(msg); resp.Send {
				resps <- resp.Text
			}
		})
	}

	go func() {
		wg.Wait()
		close(resps)
	}()

	var lines []string
	for r := range resps {
		lines = append(lines, r)
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})

	return Response{
		Text: strings.Join(lines, "\n"),
		Send: len(lines) > 0,
	}
}

// ReactOn returns combined list of all keywords
func (c Commands) ReactOn() (res []string) {
	for _, command := range c {
		res = append(res, command.ReactOn()...)
	}
	return res
}

func contains(s []string, e string) bool {
	e = strings.TrimSpace(e)
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
