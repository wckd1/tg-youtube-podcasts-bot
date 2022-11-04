package db

import "context"

const createDownload = `
INSERT INTO downloads(audio_url, cover_url, title, description)
VALUES(?,?,?,?)
RETURNING *
`

type CreateDownloadParams struct {
	AudioURL    string
	CoverURL    string
	Title       string
	Description string
}

func (q *Queries) CreateDownload(ctx context.Context, arg CreateDownloadParams) (Download, error) {
	var d Download
	stmt, err := q.db.PrepareContext(ctx, createDownload)
	if err != nil {
		return d, err
	}
	row := stmt.QueryRowContext(ctx, arg.AudioURL, arg.CoverURL, arg.Title, arg.Description)
	err = row.Scan(
		&d.ID,
		&d.AudioURL,
		&d.CoverURL,
		&d.Title,
		&d.Description,
	)
	return d, err
}
