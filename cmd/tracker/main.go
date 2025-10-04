package main

import (
	"os"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
	router.Handle("GET", "/click", handleClick)
	router.Run(":80")
}

func handleClick(c *gin.Context) {
	connectionStr := os.Getenv("DB_URL")
	if connectionStr == "" {
		connectionStr = "mysql:password@tcp(127.0.0.1:3306)/dev_db"
	}
	if connectionStr == "" {
		panic("DB_URL environment variable not set")
	}
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO clicks (account_id, campaign_id) VALUES ('1', '1')")
	if err != nil {
		panic(err.Error())
		return
	}
}
