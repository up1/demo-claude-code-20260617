package ports

import "context"

// DedupStore is a driven port for duplicate detection. IsDuplicate atomically
// records the key and reports whether it had already been seen within its TTL.
type DedupStore interface {
	IsDuplicate(ctx context.Context, key string) (bool, error)
}
