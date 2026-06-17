package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"api/internal/adapters/handler"
	"api/internal/adapters/middleware"
	"api/internal/adapters/repository"
	"api/internal/config"
	"api/internal/core/service"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		slog.Error("mongodb connect error", "error", err)
		os.Exit(1)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	if err := client.Ping(ctx, nil); err != nil {
		slog.Error("mongodb ping failed", "error", err)
		os.Exit(1)
	}

	db := client.Database(cfg.MongoDB)
	ensureIndexes(ctx, db)

	orderRepo := repository.NewOrderMongoRepo(db)
	productRepo := repository.NewProductMongoRepo(db)
	orderSvc := service.NewOrderService(orderRepo, productRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())

	api := r.Group("/api/v1")
	api.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		api.POST("/orders", orderHandler.CreateOrder)
	}

	slog.Info("server starting", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func ensureIndexes(ctx context.Context, db *mongo.Database) {
	orders := db.Collection("orders")
	_, _ = orders.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "idempotency_key", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "customer_id", Value: 1}},
		},
	})
}
