package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"api/internal/core/domain"
	"api/internal/core/ports"
)

type orderMongoRepo struct {
	col *mongo.Collection
}

func NewOrderMongoRepo(db *mongo.Database) ports.OrderRepository {
	return &orderMongoRepo{col: db.Collection("orders")}
}

func (r *orderMongoRepo) Save(ctx context.Context, order *domain.Order) error {
	_, err := r.col.InsertOne(ctx, order)
	return err
}

func (r *orderMongoRepo) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Order, error) {
	var order domain.Order
	err := r.col.FindOne(ctx, bson.M{"idempotency_key": key}).Decode(&order)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}
