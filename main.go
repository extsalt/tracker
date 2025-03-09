package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server on localhost:8000")
	http.HandleFunc("/c", clickHttpHandler)
	http.ListenAndServe("0.0.0.0:8000", nil)

}

func clickHttpHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	if !queryParams.Has("c") {
		fmt.Fprint(w, "Campaign ID is missing")
		return
	}
	if !queryParams.Has("a") {
		fmt.Fprint(w, "Affiliate ID is missing")
		return
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
	//
}
