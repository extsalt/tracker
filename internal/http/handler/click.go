package handler

import (
	"encoding/json"
	"extsalt/tracker/internal/models"
	"extsalt/tracker/internal/pubsub"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

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

	var offer models.Offer
	if err := json.Unmarshal([]byte(val), &offer); err != nil {
		panic(err)
	}

	affiliateID := c.Query("affiliate_id")
	status := "rejected"
	currentTime := time.Now().Unix()
	if slices.Contains(offer.AllowedPublishers, affiliateID) && currentTime >= offer.StartTime && currentTime <= offer.EndTime {
		status = "approved"
	}

	payload := models.ClickPayload{
		OfferID:     offerID,
		AccountID:   c.Query("account_id"),
		AffiliateID: affiliateID,
		Status:      status,
		Timestamp:   currentTime,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	redisClient.Publish(c, "click", payloadBytes)
}
