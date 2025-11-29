package main

import (
	"extsalt/tracker/internal/geo"
	"extsalt/tracker/internal/http/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := geo.Init("GeoLite2-City.mmdb"); err != nil {
		// Log error but continue as per plan
	}
	router := gin.Default()
	router.Handle("GET", "/c", handler.HandlerClick)
	router.Handle("GET", "/p", handler.HandlerConversion)
	router.Handle("GET", "/ping", handler.HandlerPing)
	err := router.Run(":80")
	if err != nil {
		panic(err)
	}
}
