package loader

import (
	"context"
	"log"
	"wckd1/tg-youtube-podcasts-bot/db"
)

type Interface interface {
	Download(src string)
}

type YTLoader struct {
	Context   context.Context
	Store     db.Store
	Submitter Submitter
}

// Submitter defines interface to submit (usually asynchronously) to the chat
type Submitter interface {
	SubmitText(ctx context.Context, text string)
	SubmitDowload(ctx context.Context, download db.Download)
}

func NewLoader(ctx context.Context, db db.Store, submitter Submitter) Interface {
	return &YTLoader{
		Context:   ctx,
		Store:     db,
		Submitter: submitter,
	}
}

func (l YTLoader) Download(src string) {
	params := db.CreateDownloadParams{
		AudioURL:    src,
		CoverURL:    "https://picsum.photos/200/300",
		Title:       "Title",
		Description: "Description",
	}
	download, err := l.Store.CreateDownload(l.Context, params)
	if err != nil {
		log.Printf("[ERROR] failed to download, %v", err)
		l.Submitter.SubmitText(l.Context, "failed to download")
		return
	}

	l.Submitter.SubmitDowload(l.Context, download)
}
