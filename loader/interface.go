package loader

import (
	"context"
	"wckd1/tg-youtube-podcasts-bot/db"
)

type Interface interface {
	Download(src string)
}

// Submitter defines interface to submit to the chat
type Submitter interface {
	SubmitText(ctx context.Context, text string)
	SubmitDowload(ctx context.Context, download db.Download)
}
