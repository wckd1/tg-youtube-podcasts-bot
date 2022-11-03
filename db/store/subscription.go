package db

import (
	"context"
	"fmt"
)

type CreateSubscriptionParams struct {
	Channel string `json:"channel"`
	Title   string `json:"title"`
}

func (q *Queries) CreateSubsctiption(ctx context.Context, arg CreateSubscriptionParams) error { // (Subscription, error)
	// TODO: Add implementation
	return fmt.Errorf("not implemented yet")
}
