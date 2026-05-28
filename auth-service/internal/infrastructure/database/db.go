package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Connect() {
	dsn := os.Getenv("AUTH_DB_URL")
	if dsn == "" {
		log.Fatal("AUTH_DB_URL is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("database: open: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("database: ping: %v", err)
	}

	DB = db
	log.Println("postgres connected")
}
