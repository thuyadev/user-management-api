package database

import (
	"fmt"
	"log"

	"user-management-api/models"
	"user-management-api/utils"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres(cfg *utils.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	logLevel := logger.Info
	if cfg.AppEnv == "production" {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return db, nil
}

func ConnectSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sqlite: %w", err)
	}
	return db, nil
}

func Connect(cfg *utils.Config) (*gorm.DB, error) {
	switch cfg.DBDriver {
	case "sqlite":
		return ConnectSQLite()
	case "postgres":
		return ConnectPostgres(cfg)
	default:
		return nil, fmt.Errorf("unsupported db driver: %s", cfg.DBDriver)
	}
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
	); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("Database migration completed")
	return nil
}
