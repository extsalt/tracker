package main

import (
	"context"
	"encoding/json"
	"extsalt/tracker/internal/db"
	"extsalt/tracker/internal/models"
	"extsalt/tracker/internal/pubsub"
	"fmt"
	"log"
)

func main() {
	redis, err := pubsub.NewClient()
	if err != nil {
		panic(err)
	}

	// Ensure DB connection is established
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

	for msg := range channel {
		tokens <- struct{}{} // Acquire token
		go func(payload string) {
			defer func() { <-tokens }() // Release token

			fmt.Println("Processing job:", payload)

			var click models.ClickPayload
			if err := json.Unmarshal([]byte(payload), &click); err != nil {
				log.Printf("Error unmarshaling payload: %v", err)
				return
			}

			_, err := database.Exec(`INSERT INTO clicks 
				(offer_id, account_id, affiliate_id, status, timestamp, ip_address, user_agent, country, state, city) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				click.OfferID, click.AccountID, click.AffiliateID, click.Status, click.Timestamp,
				click.IPAddress, click.UserAgent, click.Country, click.State, click.City)

			if err != nil {
				log.Printf("Error inserting click: %v", err)
			}
		}(msg.Payload)
	}
}
