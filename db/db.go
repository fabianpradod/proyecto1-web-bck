package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("host=localhost user=%s password= dbname=seriestracker port=5432 sslmode=disable", os.Getenv("USER"))
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	fmt.Println("Connected to PostgreSQL")
	runSchema()
}

func runSchema() {
	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatalf("Error reading schema.sql: %v", err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatalf("Error running schema: %v", err)
	}

	fmt.Println("Schema applied")
}
