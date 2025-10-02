// Package config handles the initialization and configuration of the database connection.
//
// It establishes a connection to a PostgreSQL database using GORM, configures the connection pool,
// and performs automatic migrations for the Contact model.
package config

import (
	"fmt"
	"log"
	"time"

	"api-contact-form/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DB is a global variable that holds the database connection instance.
var DB *gorm.DB

// GetEnv is assumed to exist elsewhere in your codebase. If not, uncomment this.
// func GetEnv(key, def string) string {
// 	if v := os.Getenv(key); v != "" {
// 		return v
// 	}
// 	return def
// }

// InitDB initializes the PostgreSQL connection using environment variables.
// Steps:
// 1) Read env
// 2) Build Postgres DSN (with sslmode & TimeZone suitable for local dev)
// 3) Open DB with GORM + SingularTable naming
// 4) Tune connection pool
// 5) Auto-migrate models
func InitDB() {
	_ = godotenv.Load() // ensure .env is loaded when running compiled binary

	// Retrieve configuration with safe defaults for local dev
	dbUser := GetEnv("DB_USER", "appuser")
	dbPassword := GetEnv("DB_PASSWORD", "appsecret")
	dbHost := GetEnv("DB_HOST", "127.0.0.1")
	dbPort := GetEnv("DB_PORT", "5432")
	dbName := GetEnv("DB_NAME", "contactsdb")
	sslmode := GetEnv("DB_SSLMODE", "disable")    // local dev: disable TLS
	tz := GetEnv("DB_TZ", "Asia/Jakarta")         // GORM Postgres supports TimeZone in DSN

	// DSN format per GORM Postgres driver
	// Example: host=127.0.0.1 user=appuser password=appsecret dbname=contactsdb port=5432 sslmode=disable TimeZone=Asia/Jakarta
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, sslmode, tz,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // keep your existing singular tables
		},
		// You can add Logger or other options here if needed
	})
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance!")
	}

	// Connection pool tuning (reasonable local defaults)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	// Auto-migrate your models
	if err := DB.AutoMigrate(&models.Contact{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	log.Printf("Connected to Postgres %s:%s db=%s as %s (sslmode=%s, tz=%s)",
		dbHost, dbPort, dbName, dbUser, sslmode, tz)
}
