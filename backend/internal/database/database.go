package database

import (
	"fmt"
	"log"

	"ktchat/backend/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init initializes the database connection and runs migrations
func Init(databaseURL string) (*gorm.DB, error) {
	var err error

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)

	// Connect to database
	DB, err = gorm.Open(mysql.Open(databaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")

	// Run migrations
	if err := runMigrations(DB); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return DB, nil
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Auto migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.ChatRoom{},
		&models.RoomParticipant{},
		&models.Message{},
		&models.FileAttachment{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
} 