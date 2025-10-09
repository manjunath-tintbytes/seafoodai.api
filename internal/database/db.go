package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB: ", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db
}

// RunMigrations runs all pending migrations
func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err := m.Migrate(); err != nil {
		log.Printf("Could not migrate: %v", err)
		return err
	}

	log.Printf("Migration did run successfully")
	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err := m.RollbackLast(); err != nil {
		log.Printf("Could not rollback: %v", err)
		return err
	}

	log.Printf("Rollback did run successfully")
	return nil
}

// RollbackToMigration rolls back to a specific migration
func RollbackToMigration(db *gorm.DB, migrationID string) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err := m.RollbackTo(migrationID); err != nil {
		log.Printf("Could not rollback to migration %s: %v", migrationID, err)
		return err
	}

	log.Printf("Rollback to migration %s did run successfully", migrationID)
	return nil
}
