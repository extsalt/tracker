package worker

import (
	"database/sql"
	"encoding/json"
	"extsalt/tracker/internal/models"
	"fmt"
	"log"
)

type ClickProcessor struct {
	db *sql.DB
}

func NewClickProcessor(db *sql.DB) *ClickProcessor {
	return &ClickProcessor{db: db}
}

func (p *ClickProcessor) Process(payload string) {
	fmt.Println("Processing job:", payload)

	var click models.ClickPayload
	if err := json.Unmarshal([]byte(payload), &click); err != nil {
		log.Printf("Error unmarshaling payload: %v", err)
		return
	}

	_, err := p.db.Exec(`INSERT INTO clicks 
		(offer_id, account_id, affiliate_id, status, timestamp, ip_address, user_agent, country, state, city) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		click.OfferID, click.AccountID, click.AffiliateID, click.Status, click.Timestamp,
		click.IPAddress, click.UserAgent, click.Country, click.State, click.City)

	if err != nil {
		log.Printf("Error inserting click: %v", err)
	}
}
