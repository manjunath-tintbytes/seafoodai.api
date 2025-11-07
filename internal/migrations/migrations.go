package migrations

import (
	"github.com/manjunath-tintbytes/seafoodai.api/internal/models" // Replace with your actual module path

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// GetMigrations returns all migrations
func GetMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "202409170006_create_initial_tables",
			Migrate: func(tx *gorm.DB) error {
				// Create species table
				if err := tx.AutoMigrate(&models.Species{}); err != nil {
					return err
				}

				// Create categories table
				if err := tx.AutoMigrate(&models.Category{}); err != nil {
					return err
				}

				// Create regions table
				if err := tx.AutoMigrate(&models.Region{}); err != nil {
					return err
				}

				// Create sub_regions table
				if err := tx.AutoMigrate(&models.SubRegion{}); err != nil {
					return err
				}

				// Create seafoods table
				if err := tx.AutoMigrate(&models.Seafood{}); err != nil {
					return err
				}

				// Create prices table
				if err := tx.AutoMigrate(&models.Price{}); err != nil {
					return err
				}

				// Create landing_names
				if err := tx.AutoMigrate(&models.LandingName{}); err != nil {
					return err
				}

				// Create landing_ports
				if err := tx.AutoMigrate(&models.LandingPort{}); err != nil {
					return err
				}

				// Create landings table
				if err := tx.AutoMigrate(&models.Landing{}); err != nil {
					return err
				}

				// Create market signals
				if err := tx.AutoMigrate(&models.MarketSignal{}); err != nil {
					return err
				}

				// Create quota tables
				if err := tx.AutoMigrate(&models.Quota{}); err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				// Drop tables in reverse order due to foreign key constraints
				if err := tx.Migrator().DropTable(&models.Quota{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&models.MarketSignal{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&models.Landing{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&models.LandingPort{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&models.LandingName{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&models.Price{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&models.Seafood{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&models.SubRegion{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&models.Region{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&models.Category{}); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable(&models.Species{}); err != nil {
					return err
				}
				return nil
			},
		},
	}
}
