package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // Import godotenv package
	_ "github.com/lib/pq"      // Import PostgreSQL driver
)

// InitDB untuk menginisialisasi koneksi ke database
func InitDB() (*sql.DB, error) {
	// Memuat file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Pastikan semua environment variables sudah ada
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open the database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Failed to connect to database: %v\n", err)
		return nil, err
	}

	// Tes koneksi ke database
	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v\n", err)
		return nil, err
	}
	log.Println("Database connection established successfully!")

	return db, nil
}
