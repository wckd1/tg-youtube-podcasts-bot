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

const deleteSubscription = `
DELETE FROM subscriptions
WHERE source_path = ?
`

type DeleteSubscriptionParams struct {
	SourcePath string
	Title      string
}

func (q *Queries) DeleteSubsctiption(ctx context.Context, arg DeleteSubscriptionParams) error {
	stmt, err := q.db.PrepareContext(ctx, deleteSubscription)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, arg.SourcePath)
	return err
}

const deleteTitledSubscription = `
DELETE FROM subscriptions
WHERE source_path = ?
AND title = ?
`

func (q *Queries) DeleteTitledSubsctiption(ctx context.Context, arg DeleteSubscriptionParams) error {
	stmt, err := q.db.PrepareContext(ctx, deleteTitledSubscription)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, arg.SourcePath, arg.Title)
	return err
}
