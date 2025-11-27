package handler

import (
	"extsalt/tracker/internal/pubsub"

	"github.com/gin-gonic/gin"
)


func HandlerClick(c *gin.Context) {
	redis, err := pubsub.PubSubConnect()
	if err != nil {
		panic(err)
	}
	redis.Publish(c, "click", "click")
}