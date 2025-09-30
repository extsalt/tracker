package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// http.HandleFunc("/c", clickHttpHandler)
	// http.ListenAndServe(":8000", nil)
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer redis.Close()
	// create context

	ctx := context.Background()
	err := redis.Set(ctx, "key", "value", 5*time.Minute).Err()
	if err != nil {
		panic(err)
	}
	val, err := redis.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}

func FullURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host + r.RequestURI
}

func clickHttpHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	if !queryParams.Has("c") {
		fmt.Fprint(w, "Campaign ID is missing")
	}
	if !queryParams.Has("a") {
		fmt.Fprint(w, "Affiliate ID is missing")
	}
	// connect with redis and check if campain exists
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer redis.Close()
	err := redis.Set(r.Context(), "key", "value", 5*time.Minute).Err()
	if err != nil {
		panic(err)
	}

	// Get campaigns details from redis and update the click count
	// Redirect to the campaign URL
	// Check 1 - Check if the campaign exists
	// Check 2 - Check if the affiliate exists
	// Check 3 - Check if the campaign is active
	// Check 4 - Check if the affiliate is active
	// Check 5 - Check if the campaign is active for the affiliate
	// Check 6 - Check if the campaign is active for the affiliate for the current date
	// Check 7 - Check if the campaign is active for the affiliate for the current date and time
	// Check 8 - Check if the campaign is active for the affiliate for the current date and time
	// get full path from the request
	fmt.Println("Full path:", FullURL(r))
}
