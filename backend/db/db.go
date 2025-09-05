package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := "postgres://postgres:1234@127.0.0.1:5432/kanban?sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connect error:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}
	fmt.Println("Connected to PostgreSQL")
}
