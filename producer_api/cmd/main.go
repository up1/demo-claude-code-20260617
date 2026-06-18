package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"producer/internal/adapters/cache"
	"producer/internal/adapters/handler"
	"producer/internal/adapters/middleware"
	"producer/internal/adapters/publisher"
	"producer/internal/config"
	"producer/internal/core/service"
	"producer/internal/observability"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "error", err)
		os.Exit(1)
	}

	rootCtx := context.Background()

	// OpenTelemetry tracing.
	shutdownTracer, err := observability.InitTracer(rootCtx, cfg.ServiceName, cfg.OTLPEndpoint)
	if err != nil {
		slog.Error("otel init error", "error", err)
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTracer(ctx)
	}()

	// Redis (dedup store).
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	defer func() { _ = redisClient.Close() }()

	pingCtx, cancelPing := context.WithTimeout(rootCtx, 10*time.Second)
	defer cancelPing()
	if err := redisClient.Ping(pingCtx).Err(); err != nil {
		slog.Error("redis ping failed", "error", err, "addr", cfg.RedisAddr)
		os.Exit(1)
	}

	// RabbitMQ (publisher).
	connCtx, cancelConn := context.WithTimeout(rootCtx, 30*time.Second)
	defer cancelConn()
	rmqPublisher, err := publisher.NewRabbitMQPublisher(connCtx, cfg.RabbitMQURI, service.QueueName)
	if err != nil {
		slog.Error("rabbitmq init error", "error", err)
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = rmqPublisher.Close(ctx)
	}()

	// Wire the use-case and HTTP adapter.
	dedup := cache.NewRedisDedup(redisClient, cfg.DedupTTL)
	lineSvc := service.NewLineService(rmqPublisher, dedup)
	lineHandler := handler.NewLineHandler(lineSvc)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(otelgin.Middleware(cfg.ServiceName))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/line/messages", lineHandler.ReceiveMessage)
	}

	slog.Info("server starting", "port", cfg.Port, "service", cfg.ServiceName)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
