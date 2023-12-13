package file_manager

import (
	"context"
	"time"
)

type FileManager interface {
	Get(ctx context.Context, url string) (download Download, err error)
	CheckUpdate(ctx context.Context, url string, date time.Time, filter string) (downloads []Download, err error)
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
}

type Download struct {
	URL  string
	Info FileInfo
}
