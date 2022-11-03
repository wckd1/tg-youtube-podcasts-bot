package db

import (
	"context"
)

const createSubscription = `
INSERT INTO subscriptions(channel, title)
VALUES(?,?)
`

type CreateSubscriptionParams struct {
	Channel string `json:"channel"`
	Title   string `json:"title"`
}

func (q *Queries) CreateSubsctiption(ctx context.Context, arg CreateSubscriptionParams) error {
	stmt, err := q.db.PrepareContext(ctx, createSubscription)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, arg.Channel, arg.Title)
	return err
}
