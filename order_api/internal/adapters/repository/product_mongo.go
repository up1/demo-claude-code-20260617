package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"api/internal/core/domain"
	"api/internal/core/ports"
)

type productMongoRepo struct {
	col *mongo.Collection
}

func NewProductMongoRepo(db *mongo.Database) ports.ProductRepository {
	return &productMongoRepo{col: db.Collection("products")}
}

func (r *productMongoRepo) FindByID(ctx context.Context, productID string) (*domain.Product, error) {
	var product domain.Product
	err := r.col.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}
