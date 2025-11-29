package models

type Offer struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	AllowedPublishers []string `json:"allowed_publishers"`
	StartTime         int64    `json:"start_time"`
	EndTime           int64    `json:"end_time"`
	OfferURL          string   `json:"offer_url"`
	FallbackURL       string   `json:"fallback_url"`
}

type ClickPayload struct {
	OfferID     string `json:"offer_id"`
	AccountID   string `json:"account_id"`
	AffiliateID string `json:"affiliate_id"`
	Status      string `json:"status"`
	Timestamp   int64  `json:"timestamp"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
	Country     string `json:"country"`
	State       string `json:"state"`
	City        string `json:"city"`
}
