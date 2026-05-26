package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgre() (*gorm.DB, error) {
	// 1. Gather variables with production-ready default fallbacks
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password123"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "emission_preferences"
	}
	schema := os.Getenv("DB_SCHEMA")
	if schema == "" {
		schema = "public"
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable" // Switch to 'require' or 'verify-full' for remote cloud dbs
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s search_path=%s",
		host, user, password, dbname, port, sslmode, schema,
	)

	// 2. Setup database logger to match standard application logs
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Log slow SQL executions
			LogLevel:                  logger.Warn,            // Error & Warn logs only in prod
			IgnoreRecordNotFoundError: true,                   // Keeps log clean from standard 404 queries
			Colorful:                  true,
		},
	)

	// 3. Establish connection
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("notification-service: postgres connect error: ", err)
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("notification-service: postgres sql-instance recovery error: ", err)
		return nil, err
	}

	// 4. Configure connection pool bounds
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Minute * 10)

	// 5. Verify connection using a synchronous Ping
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("notification-service: postgres ping error: ", err)
		return nil, err
	}

	log.Printf("notification-service: connected to postgres database: %s (schema: %s)", dbname, schema)
	return gormDB, nil
}
