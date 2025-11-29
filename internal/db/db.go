package db

import (
	"database/sql"
	"os"
)

func DBConnect() (*sql.DB, error) {
	connectionStr := os.Getenv("DB_URL")
	if connectionStr == "" {
		connectionStr = "mysql:password@tcp(127.0.0.1:3306)/dev_db?charset=utf8mb4&parseTime=True&loc=Local"
	}
	if connectionStr == "" {
		panic("DB_URL environment variable not set")
	}
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		panic(err.Error())
	}
	return db, err
}
