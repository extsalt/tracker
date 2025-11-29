package main

import (
	"context"
	"extsalt/tracker/internal/db"
	"extsalt/tracker/internal/pubsub"
	"extsalt/tracker/internal/worker"
)

func main() {
	redis, err := pubsub.NewClient()
	if err != nil {
		panic(err)
	}
	database, err := db.DBConnect()
	if err != nil {
		panic(err)
	}

	// Create table if not exists
	_, err = database.Exec(`CREATE TABLE IF NOT EXISTS clicks (
		id INT AUTO_INCREMENT PRIMARY KEY,
		offer_id VARCHAR(255),
		account_id VARCHAR(255),
		affiliate_id VARCHAR(255),
		status VARCHAR(50),
		timestamp BIGINT,
		ip_address VARCHAR(45),
		user_agent TEXT,
		country VARCHAR(2),
		state VARCHAR(10),
		city VARCHAR(255)
	)`)
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
	clickProcessor := worker.NewClickProcessor(database)

	for msg := range channel {
		tokens <- struct{}{} // Acquire token
		go func(payload string) {
			defer func() { <-tokens }() // Release token
			clickProcessor.Process(payload)
		}(msg.Payload)
	}
}
