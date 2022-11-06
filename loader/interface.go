package loader

import (
	"context"
	"wckd1/tg-youtube-podcasts-bot/db"
)

// TODO: Separeate in two services
type Interface interface {
	Download(src string)
	Upload(download db.Download)
}

// Submitter defines interface to submit to the chat
type Submitter interface {
	SubmitText(ctx context.Context, text string)
}
