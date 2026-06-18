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

	inboxRepo := repository.NewInboxMongoRepo(db)
	inboxSvc := service.NewInboxService(inboxRepo)
	inboxHandler := handler.NewInboxHandler(inboxSvc)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())

	api := r.Group("/api/v1")
	api.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		api.GET("/inbox/messages", inboxHandler.ListMessages)
	}

	slog.Info("server starting", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func ensureIndexes(ctx context.Context, db *mongo.Database) {
	inbox := db.Collection("inbox")
	_, _ = inbox.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			// Default sort.
			Keys: bson.D{{Key: "updated_at", Value: -1}},
		},
		{
			// Filtered list + sort.
			Keys: bson.D{
				{Key: "channel", Value: 1},
				{Key: "status", Value: 1},
				{Key: "updated_at", Value: -1},
			},
		},
	})
}
