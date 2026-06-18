package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcmongo "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"api/internal/adapters/repository"
	"api/internal/core/domain"
	"api/internal/core/ports"
)

func setupMongo(t *testing.T) *mongo.Database {
	t.Helper()
	ctx := context.Background()

	container, err := tcmongo.Run(ctx, "mongo:8")
	require.NoError(t, err)
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })

	uri, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Disconnect(ctx) })

	require.NoError(t, client.Ping(ctx, nil))
	return client.Database("inbox_test")
}

func seed(t *testing.T, db *mongo.Database) {
	t.Helper()
	base := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	docs := []interface{}{
		domain.InboxMessage{CustomerID: "c1", SenderName: "Marcus Watanabe", Channel: domain.ChannelLine, Preview: "Can you confirm the tracking number?", Status: domain.StatusPending, Unread: true, CreatedAt: base, UpdatedAt: base.Add(4 * time.Hour)},
		domain.InboxMessage{CustomerID: "c2", SenderName: "Sarah Jenkins", Channel: domain.ChannelFacebook, Preview: "Thank you!", Status: domain.StatusReplied, Unread: false, CreatedAt: base, UpdatedAt: base.Add(3 * time.Hour)},
		domain.InboxMessage{CustomerID: "c3", SenderName: "Aiko Tanaka", Channel: domain.ChannelLine, Preview: "Where is my refund?", Status: domain.StatusPending, Unread: true, CreatedAt: base, UpdatedAt: base.Add(2 * time.Hour)},
		domain.InboxMessage{CustomerID: "c4", SenderName: "John Smith", Channel: domain.ChannelInstagram, Preview: "Tracking link broken", Status: domain.StatusReplied, Unread: false, CreatedAt: base, UpdatedAt: base.Add(1 * time.Hour)},
	}
	_, err := db.Collection("inbox").InsertMany(context.Background(), docs)
	require.NoError(t, err)
}

func TestInboxMongo_List(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testcontainers integration test in -short mode")
	}

	db := setupMongo(t)
	seed(t, db)
	repo := repository.NewInboxMongoRepo(db)
	ctx := context.Background()

	t.Run("default sort by updated_at desc", func(t *testing.T) {
		items, total, err := repo.List(ctx, ports.InboxFilter{Offset: 0, Limit: 20})
		require.NoError(t, err)
		assert.Equal(t, int64(4), total)
		require.Len(t, items, 4)
		for i := 1; i < len(items); i++ {
			assert.False(t, items[i-1].UpdatedAt.Before(items[i].UpdatedAt), "items must be sorted by updated_at desc")
		}
		assert.Equal(t, "Marcus Watanabe", items[0].SenderName)
	})

	t.Run("filter by channel", func(t *testing.T) {
		items, total, err := repo.List(ctx, ports.InboxFilter{Channel: "line", Offset: 0, Limit: 20})
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		for _, it := range items {
			assert.Equal(t, domain.ChannelLine, it.Channel)
		}
	})

	t.Run("filter by status", func(t *testing.T) {
		_, total, err := repo.List(ctx, ports.InboxFilter{Status: "pending", Offset: 0, Limit: 20})
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
	})

	t.Run("combined filter and search", func(t *testing.T) {
		items, total, err := repo.List(ctx, ports.InboxFilter{Channel: "line", Search: "tracking", Offset: 0, Limit: 20})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		require.Len(t, items, 1)
		assert.Equal(t, "Marcus Watanabe", items[0].SenderName)
	})

	t.Run("pagination skip and limit", func(t *testing.T) {
		page1, total, err := repo.List(ctx, ports.InboxFilter{Offset: 0, Limit: 2})
		require.NoError(t, err)
		assert.Equal(t, int64(4), total)
		require.Len(t, page1, 2)

		page2, _, err := repo.List(ctx, ports.InboxFilter{Offset: 2, Limit: 2})
		require.NoError(t, err)
		require.Len(t, page2, 2)
		assert.NotEqual(t, page1[0].CustomerID, page2[0].CustomerID)
	})

	t.Run("empty result", func(t *testing.T) {
		items, total, err := repo.List(ctx, ports.InboxFilter{Search: "no-such-text-anywhere", Offset: 0, Limit: 20})
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, items)
		assert.NotNil(t, items)
	})
}
