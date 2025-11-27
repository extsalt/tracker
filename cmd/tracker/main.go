package main

import (
	"extsalt/tracker/internal/http/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
	router.Handle("GET", "/c", handler.HandlerClick)
	router.Handle("GET", "/p", handler.HandlerConversion)
	router.Handle("GET", "/ping", handler.HandlerPing)
	err := router.Run(":80")
	if err != nil {
		panic(err)
	}
}
