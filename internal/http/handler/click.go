package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"slices"
	"time"

	"extsalt/tracker/internal/geo"
	"extsalt/tracker/internal/models"
	"extsalt/tracker/internal/pubsub"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func HandlerClick(c *gin.Context) {
	offerID := c.Query("offer_id")
	if offerID == "" {
		c.JSON(400, gin.H{"error": "offer_id is required"})
		return
	}
	affiliateID := c.Query("affiliate_id")
	if affiliateID == "" {
		c.JSON(400, gin.H{"error": "affiliate_id is required"})
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

	status := "approved"
	currentTime := time.Now().Unix()
	if !slices.Contains(offer.AllowedPublishers, affiliateID) {
		status = "rejected"
	}
	if currentTime < offer.StartTime || currentTime > offer.EndTime {
		status = "rejected"
	}

	country, state, city, _ := geo.Lookup(c.ClientIP())

	payload := models.ClickPayload{
		OfferID:     offerID,
		AccountID:   c.Query("account_id"),
		AffiliateID: affiliateID,
		Status:      status,
		Timestamp:   currentTime,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		Country:     country,
		State:       state,
		City:        city,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	redisClient.Publish(c, "click", payloadBytes)

	targetURL := offer.OfferURL
	if status == "rejected" && offer.FallbackURL != "" {
		if _, err := url.ParseRequestURI(offer.FallbackURL); err == nil {
			targetURL = offer.FallbackURL
		}
	}

	c.Redirect(http.StatusFound, targetURL)
}
