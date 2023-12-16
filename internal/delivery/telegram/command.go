package telegram

import (
	"sort"
	"strings"
	"sync"
)

type Command interface {
	OnMessage(msg Message) Response
	ReactOn() (res []string)
}

// Commands combines many commands to one virtual
type commandList []Command

// OnMessage pass msg to all commands and collects responses (combining all of them)
func (cl commandList) OnMessage(msg Message) Response {
	resps := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(len(cl))

	for _, c := range cl {
		go func(cmd Command) {
			defer wg.Done()

			if resp := cmd.OnMessage(msg); resp.Send {
				resps <- resp.Text
			}
		}(c)
	}

	wg.Wait()
	close(resps)

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
func (cl commandList) ReactOn() (res []string) {
	for _, command := range cl {
		res = append(res, command.ReactOn()...)
	}
	return res
}
