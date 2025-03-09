package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/c", clickHttpHandler)
	http.ListenAndServe("localhost:8000", nil)
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

}
