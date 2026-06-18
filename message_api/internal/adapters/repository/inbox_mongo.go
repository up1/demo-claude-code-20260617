package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"api/internal/core/domain"
	"api/internal/core/ports"
)

type inboxMongoRepo struct {
	col *mongo.Collection
}

// NewInboxMongoRepo returns a MongoDB-backed InboxRepository on the "inbox" collection.
func NewInboxMongoRepo(db *mongo.Database) ports.InboxRepository {
	return &inboxMongoRepo{col: db.Collection("inbox")}
}

// List builds the query (filters + case-insensitive regex search), counts the
// total matching documents, then returns the requested page sorted by
// updated_at descending.
func (r *inboxMongoRepo) List(ctx context.Context, filter ports.InboxFilter) ([]domain.InboxMessage, int64, error) {
	query := buildQuery(filter)

	total, err := r.col.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: -1}}).
		SetSkip(int64(filter.Offset)).
		SetLimit(int64(filter.Limit))

	cur, err := r.col.Find(ctx, query, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = cur.Close(ctx) }()

	messages := make([]domain.InboxMessage, 0)
	if err := cur.All(ctx, &messages); err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func buildQuery(filter ports.InboxFilter) bson.M {
	query := bson.M{}

	if filter.Channel != "" {
		query["channel"] = filter.Channel
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.From != nil || filter.To != nil {
		updatedAt := bson.M{}
		if filter.From != nil {
			updatedAt["$gte"] = *filter.From
		}
		if filter.To != nil {
			updatedAt["$lte"] = *filter.To
		}
		query["updated_at"] = updatedAt
	}

	if filter.Search != "" {
		rgx := bson.M{"$regex": filter.Search, "$options": "i"}
		query["$or"] = bson.A{
			bson.M{"sender_name": rgx},
			bson.M{"preview": rgx},
		}
	}

	return query
}
