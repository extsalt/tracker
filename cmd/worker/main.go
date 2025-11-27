package main

import (
	"context"
	"extsalt/tracker/internal/pubsub"
	"fmt"
)

func main() {
	redis, err := pubsub.PubSubConnect()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	subscription := redis.Subscribe(ctx, "click")
	_, err1 := subscription.Receive(ctx)
	if err1 != nil {
		panic(err1)
	}
	channel  := subscription.Channel()
	for msg := range channel {
		fmt.Println(msg.Payload)
	}
}
