package content

import (
	"context"
	"time"
)

type ContentManager interface {
	Get(ctx context.Context, url string) (dl Download, err error)
	CheckUpdate(ctx context.Context, url string, date time.Time, filter string) (dls []Download, err error)
}

type FileInfo struct {
	Link        string `json:"webpage_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"thumbnail"`
	Author      string `json:"uploader"`
	Length      int    `json:"filesize"`
	Duration    int    `json:"duration"`
	Date        string `json:"upload_date"`
	Type        string
}

type Download struct {
	URL  string
	Info FileInfo
}
