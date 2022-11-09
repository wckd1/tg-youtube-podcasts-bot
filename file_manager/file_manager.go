package file_manager

import (
	"context"
)

type localFile struct {
	path string
	info fileInfo
}

type fileInfo struct {
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
	Info fileInfo
}

// Downloader defines interface to download file to local fs
type Downloader interface {
	Download(ctx context.Context, id string) (file localFile, err error)
}

// Uloader defines interface to upload file from local fs to storage
type Uploader interface {
	Upload(ctx context.Context, file localFile) (url string, err error)
}

type FileManager struct {
	Downloader Downloader
	Uploader   Uploader
}

func (fm FileManager) Get(ctx context.Context, id string) (download Download, err error) {
	download = Download{}

	file, err := fm.Downloader.Download(ctx, id)
	if err != nil {
		return
	}

	download.Info = file.info

	url, err := fm.Uploader.Upload(ctx, file)
	if err != nil {
		return
	}

	download.URL = url
	return
}
