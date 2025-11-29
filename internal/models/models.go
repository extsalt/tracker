package models

type Offer struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	AllowedPublishers []string `json:"allowed_publishers"`
	StartTime         int64    `json:"start_time"`
	EndTime           int64    `json:"end_time"`
}

type ClickPayload struct {
	OfferID     string `json:"offer_id"`
	AccountID   string `json:"account_id"`
	AffiliateID string `json:"affiliate_id"`
	Status      string `json:"status"`
	Timestamp   int64  `json:"timestamp"`
}
