package handler

import (
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
	redisClient, err := pubsub.NewClient()
	if err != nil {
		panic(err)
	}

	_, err = redisClient.Get(c, offerID).Result()
	if err == redis.Nil {
		c.JSON(404, gin.H{"error": "offer not found"})
		return
	} else if err != nil {
		panic(err)
	}

	redisClient.Publish(c, "click", "click")
}
