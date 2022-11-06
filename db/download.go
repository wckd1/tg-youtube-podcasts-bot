package db

import "context"

const createDownload = `
INSERT INTO downloads(path, cover_url, title, description)
VALUES(?,?,?,?)
RETURNING *
`

type CreateDownloadParams struct {
	Path        string
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
	row := stmt.QueryRowContext(ctx, arg.Path, arg.CoverURL, arg.Title, arg.Description)
	err = row.Scan(
		&d.ID,
		&d.Path,
		&d.CoverURL,
		&d.Title,
		&d.Description,
	)
	return d, err
}
