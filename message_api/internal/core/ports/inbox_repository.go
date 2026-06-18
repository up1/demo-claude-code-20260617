package ports

import (
	"context"
	"time"

	"api/internal/core/domain"
)

// InboxFilter is the normalized query passed from the service to the repository.
// Empty strings mean "no filter"; nil time pointers mean "unbounded".
type InboxFilter struct {
	Channel string     // "" = all
	Status  string     // "" = all
	Search  string     // "" = none
	From    *time.Time // nil = unbounded
	To      *time.Time // nil = unbounded
	Offset  int        // (page-1) * pageSize
	Limit   int        // pageSize
}

// InboxRepository is the driven port implemented by the MongoDB adapter.
// List returns the requested page plus the total count of matching documents.
type InboxRepository interface {
	List(ctx context.Context, filter InboxFilter) ([]domain.InboxMessage, int64, error)
}
