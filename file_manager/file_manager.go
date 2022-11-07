package file_manager

import (
	"context"
)

type localFile struct {
	path string
	info fileInfo
}

type fileInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"thumbnail"`
}

type Download struct {
	URL         string
	Title       string
	Description string
	CoverURL    string
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

	download.Title = file.info.Title
	download.Description = file.info.Description
	download.CoverURL = file.info.ImageURL

	url, err := fm.Uploader.Upload(ctx, file)
	if err != nil {
		return
	}

	download.URL = url
	return
}
