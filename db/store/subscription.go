package db

import (
	"context"
)

const createSubscription = `
INSERT INTO subscriptions(source_path, source_type, title)
VALUES(?,?,?)
`

type CreateSubscriptionParams struct {
	SourcePath string
	SourceType SourceType
	Title      string
}

func (q *Queries) CreateSubsctiption(ctx context.Context, arg CreateSubscriptionParams) error {
	stmt, err := q.db.PrepareContext(ctx, createSubscription)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, arg.SourcePath, arg.SourceType, arg.Title)
	return err
}
