package database

import (
	"database/sql"
	"fmt"
	"time"

	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

// Connect connects to database and returns DB
func Connect(dsn string) (*DB, error) {
	// Charmbracelet's "log" as a slog handler
	charmHandler := log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		Prefix:          "GORM",
	})

	slogger := slog.New(charmHandler)

	// Create GORM logger with slog
	gormLogger := logger.New(
		slog.NewLogLogger(slogger.Handler(), slog.LevelDebug),
		logger.Config{
			SlowThreshold:             time.Millisecond * 200,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Maximum open connections
	sqlDB.SetMaxOpenConns(25)

	// Maximum idle connections
	sqlDB.SetMaxIdleConns(5)

	// Maximum lifetime for connections
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Maximum idle time for connections
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Successfully connected to database")

	return &DB{DB: db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetStats returns database connection stats
func (db *DB) GetStats() sql.DBStats {
	sqlDB, _ := db.DB.DB()
	return sqlDB.Stats()
}
