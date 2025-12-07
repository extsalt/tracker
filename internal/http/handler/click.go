package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"slices"
	"strings"
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

	ua := c.Request.UserAgent()

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

	if len(offer.AllowedUserAgents) > 0 {
		uaAllowed := false
		for _, allowed := range offer.AllowedUserAgents {
			if strings.Contains(ua, allowed) {
				uaAllowed = true
				break
			}
		}
		if !uaAllowed {
			status = "rejected"
		}
	}
	if currentTime < offer.StartTime || currentTime > offer.EndTime {
		status = "rejected"
	}

	targetURL := offer.OfferURL
	var fallbackURL string
	isFallback := false

	// Check for Affiliate Priority Fallback
	if setting, exists := offer.AffiliateSettings[affiliateID]; exists {
		if setting.EnableFallback && setting.FallbackURL != "" {
			fallbackURL = setting.FallbackURL
		}
	}

	// If no Affiliate Fallback, check for Offer Fallback
	if fallbackURL == "" && offer.EnableFallback && offer.FallbackURL != "" {
		fallbackURL = offer.FallbackURL
	}

	if status == "rejected" && fallbackURL != "" {
		if _, err := url.ParseRequestURI(fallbackURL); err == nil {
			targetURL = fallbackURL
			isFallback = true
		}
	}

	country, state, city, _ := geo.Lookup(c.ClientIP())

	payload := models.ClickPayload{
		OfferID:     offerID,
		AccountID:   c.Query("account_id"),
		AffiliateID: affiliateID,
		Status:      status,
		Timestamp:   currentTime,
		IPAddress:   c.ClientIP(),
		UserAgent:   ua,
		Country:     country,
		State:       state,
		City:        city,
		IsFallback:  isFallback,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	redisClient.Publish(c, "click", payloadBytes)

	c.Redirect(http.StatusFound, targetURL)
}
