// Command seed bulk-inserts mock InboxMessage documents into MongoDB for load/testing.
//
// Configuration (env vars):
//
//	MONGO_URI     MongoDB connection string         (default mongodb://localhost:27017)
//	MONGO_DB      Database name                     (default inbox_db)
//	SEED_COUNT    Total documents to insert         (default 10000000)
//	SEED_BATCH    Documents per InsertMany batch    (default 10000)
//	SEED_WORKERS  Concurrent insert workers         (default 8)
//	SEED_DROP     Drop the collection first (true)  (default true)
package main

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"api/internal/core/domain"
)

var (
	channels   = []domain.Channel{domain.ChannelFacebook, domain.ChannelLine, domain.ChannelInstagram}
	statuses   = []domain.Status{domain.StatusPending, domain.StatusReplied}
	firstNames = []string{"Marcus", "Sarah", "Aiko", "John", "Mei", "Carlos", "Priya", "Liam", "Yuki", "Fatima", "Noah", "Sofia", "Kenji", "Olivia", "Ahmed", "Emma", "Wei", "Lucas", "Ananya", "Diego"}
	lastNames  = []string{"Watanabe", "Jenkins", "Tanaka", "Smith", "Chen", "Garcia", "Patel", "Murphy", "Sato", "Khan", "Brown", "Rossi", "Kim", "Nguyen", "Ali", "Johnson", "Zhang", "Silva", "Sharma", "Lopez"}
	previews   = []string{
		"Can you confirm the tracking number for order #%d?",
		"Thank you for the quick resolution!",
		"Where is my refund for transaction %d?",
		"The tracking link seems broken, please help.",
		"I'd like to change my delivery address.",
		"Is this item still in stock?",
		"My package #%d arrived damaged.",
		"Can I get an invoice for my last purchase?",
		"How long does shipping usually take?",
		"Please cancel my subscription.",
	}
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	uri := envDefault("MONGO_URI", "mongodb://localhost:27017")
	dbName := envDefault("MONGO_DB", "inbox_db")
	count := envInt("SEED_COUNT", 10_000_000)
	batchSize := envInt("SEED_BATCH", 10_000)
	workers := envInt("SEED_WORKERS", 8)
	drop := envDefault("SEED_DROP", "true") == "true"

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(uint64(workers)+2))
	if err != nil {
		slog.Error("mongodb connect error", "error", err)
		os.Exit(1)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	if err := client.Ping(ctx, nil); err != nil {
		slog.Error("mongodb ping failed", "error", err)
		os.Exit(1)
	}

	col := client.Database(dbName).Collection("inbox")
	if drop {
		if err := col.Drop(ctx); err != nil {
			slog.Error("drop collection failed", "error", err)
			os.Exit(1)
		}
		slog.Info("dropped existing collection", "db", dbName, "collection", "inbox")
	}

	slog.Info("seeding started", "count", count, "batch", batchSize, "workers", workers)
	start := time.Now()

	// Each batch is one job described by its size; workers generate + insert.
	jobs := make(chan int, workers*2)
	var inserted atomic.Int64
	var wg sync.WaitGroup

	insertOpts := options.InsertMany().SetOrdered(false)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(seed))
			for n := range jobs {
				docs := make([]interface{}, n)
				for i := 0; i < n; i++ {
					docs[i] = randomMessage(rng)
				}
				if _, err := col.InsertMany(ctx, docs, insertOpts); err != nil {
					slog.Error("insert batch failed", "error", err)
					os.Exit(1)
				}
				done := inserted.Add(int64(n))
				if done%(int64(batchSize)*20) == 0 {
					rate := float64(done) / time.Since(start).Seconds()
					slog.Info("progress", "inserted", done, "total", count, "docs_per_sec", int64(rate))
				}
			}
		}(int64(w)*7919 + time.Now().UnixNano())
	}

	for remaining := count; remaining > 0; remaining -= batchSize {
		n := batchSize
		if remaining < batchSize {
			n = remaining
		}
		jobs <- n
	}
	close(jobs)
	wg.Wait()

	slog.Info("seeding complete",
		"inserted", inserted.Load(),
		"elapsed", time.Since(start).String(),
		"docs_per_sec", int64(float64(inserted.Load())/time.Since(start).Seconds()),
	)
}

func randomMessage(rng *rand.Rand) domain.InboxMessage {
	first := firstNames[rng.Intn(len(firstNames))]
	last := lastNames[rng.Intn(len(lastNames))]

	// created_at within the last 180 days; updated_at between created_at and now.
	now := time.Now().UTC()
	createdOffset := time.Duration(rng.Int63n(int64(180 * 24 * time.Hour)))
	createdAt := now.Add(-createdOffset)
	updatedAt := createdAt.Add(time.Duration(rng.Int63n(int64(createdOffset) + 1)))

	return domain.InboxMessage{
		CustomerID: "cust_" + strconv.Itoa(rng.Intn(2_000_000)),
		SenderName: first + " " + last,
		AvatarURL:  "https://i.pravatar.cc/150?u=" + first + last + strconv.Itoa(rng.Intn(100000)),
		Channel:    channels[rng.Intn(len(channels))],
		Preview:    fmtPreview(rng),
		Status:     statuses[rng.Intn(len(statuses))],
		Unread:     rng.Intn(2) == 0,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func fmtPreview(rng *rand.Rand) string {
	tmpl := previews[rng.Intn(len(previews))]
	// Templates with a %d get a random order/transaction number; others are returned as-is.
	if containsVerb(tmpl) {
		return sprintfInt(tmpl, rng.Intn(90000)+10000)
	}
	return tmpl
}

func containsVerb(s string) bool {
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '%' && s[i+1] == 'd' {
			return true
		}
	}
	return false
}

func sprintfInt(tmpl string, n int) string {
	// tiny wrapper to avoid importing fmt just for one call site
	out := make([]byte, 0, len(tmpl)+8)
	num := strconv.Itoa(n)
	for i := 0; i < len(tmpl); i++ {
		if i+1 < len(tmpl) && tmpl[i] == '%' && tmpl[i+1] == 'd' {
			out = append(out, num...)
			i++
			continue
		}
		out = append(out, tmpl[i])
	}
	return string(out)
}

func envDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}
