package main

import (
	"context"
	"extsalt/tracker/internal/db"
	"extsalt/tracker/internal/pubsub"
	"fmt"
	"time"
)

func main() {
	redis, err := pubsub.NewClient()
	if err != nil {
		panic(err)
	}

	// Ensure DB connection is established
	_, err = db.DBConnect()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	subscription := redis.Subscribe(ctx, "click")
	_, err1 := subscription.Receive(ctx)
	if err1 != nil {
		panic(err1)
	}
	channel := subscription.Channel()

	// Token bucket for concurrency control (e.g., 10 concurrent workers)
	tokens := make(chan struct{}, 10)

	for msg := range channel {
		tokens <- struct{}{} // Acquire token
		go func(payload string) {
			defer func() { <-tokens }() // Release token

			fmt.Println("Processing job:", payload)
			// Simulate work
			time.Sleep(100 * time.Millisecond)
			// Future: Add actual DB processing logic here
		}(msg.Payload)
	}
}
