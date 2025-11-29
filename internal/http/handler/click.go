package handler

import (
	"encoding/json"
	"extsalt/tracker/internal/pubsub"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Offer struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	AllowedPublishers []string `json:"allowed_publishers"`
}

type ClickPayload struct {
	OfferID     string `json:"offer_id"`
	AccountID   string `json:"account_id"`
	AffiliateID string `json:"affiliate_id"`
	Status      string `json:"status"`
}

func HandlerClick(c *gin.Context) {
	offerID := c.Query("offer_id")
	if offerID == "" {
		c.JSON(400, gin.H{"error": "offer_id is required"})
		return
	}
	redisClient, err := pubsub.NewClient()
	if err != nil {
		panic(err)
	}

	val, err := redisClient.Get(c, offerID).Result()
	if err == redis.Nil {
		c.JSON(404, gin.H{"error": "offer not found"})
		return
	} else if err != nil {
		panic(err)
	}

	var offer Offer
	if err := json.Unmarshal([]byte(val), &offer); err != nil {
		panic(err)
	}

	affiliateID := c.Query("affiliate_id")
	status := "rejected"
	if slices.Contains(offer.AllowedPublishers, affiliateID) {
		status = "approved"
	}

	payload := ClickPayload{
		OfferID:     offerID,
		AccountID:   c.Query("account_id"),
		AffiliateID: affiliateID,
		Status:      status,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	redisClient.Publish(c, "click", payloadBytes)
}
