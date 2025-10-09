package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/config"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/database"
)

/// ---------- MODELS ---------- ///

// Species table
type Species struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Categories table
type Category struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Regions table
type Region struct {
	ID        uint   `gorm:"primaryKey"`
	Region    string `gorm:"unique"`
	Quota     float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// SubRegions table (optional – we’ll skip unless you want to map SizeRange)
type SubRegion struct {
	ID        uint `gorm:"primaryKey"`
	RegionID  uint
	SubRegion string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Transactions table
type Seafood struct {
	ID          uint `gorm:"primaryKey"`
	SpeciesID   uint
	RegionID    uint
	SubRegionID *uint
	CategoryID  uint
	VolumeMt    *float64
	PriceUnit   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// Prices table
type Price struct {
	ID        uint `gorm:"primaryKey"`
	SeafoodID uint
	Date      time.Time
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

/// ---------- CSV STRUCT ---------- ///

type PriceCSV struct {
	Year      int      `csv:"YEAR"`
	Month     int      `csv:"MONTH"`
	Day       int      `csv:"DAY"`
	Country   string   `csv:"COUNTRY"`
	Category  string   `csv:"CATEGORY"`
	SizeRange string   `csv:"SIZE/WEIGHT RANGE"`
	Product   string   `csv:"PRODUCT"`
	Item      string   `csv:"ITEM ON SALE"`
	PriceUnit *float64 `csv:"PRICE PER UNIT (EUR)"`
	PriceKg   *float64 `csv:"PRICE PER KG (EUR)"`
}

/// ---------- MAIN SEEDER ---------- ///

func main() {
	config.LoadEnv()
	db := database.SetupDB()

	// Load CSV
	records, err := LoadCSV("2025_Daily-online-retail-prices.csv")
	if err != nil {
		log.Fatal("❌ failed to load CSV:", err)
	}

	// Before the loop
	seen := make(map[string]bool)

	for _, r := range records {
		// Parse date
		dateStr := fmt.Sprintf("%04d-%02d-%02d", r.Year, r.Month, r.Day)
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Println("⚠️ skipping row due to invalid date:", dateStr)
			continue
		}

		// Decide price + unit
		priceUnit := ""
		var priceValue float64
		if r.PriceKg != nil && *r.PriceKg > 0 {
			priceUnit = "kg"
			priceValue = *r.PriceKg
		} else if r.PriceUnit != nil && *r.PriceUnit > 0 {
			priceUnit = "unit"
			priceValue = *r.PriceUnit
		} else {
			log.Println("⚠️ no price found, skipping transaction:", r.Product)
			continue
		}

		// Create a unique key for deduplication
		key := fmt.Sprintf("%s|%s|%s|%s|%s|%.2f|%s",
			r.Product, r.Category, r.Country, priceUnit, dateStr, priceValue, r.SizeRange)

		if seen[key] {
			// Skip duplicate row
			continue
		}
		seen[key] = true

		// --- Insert Species ---
		var species Species
		db.Where("name = ?", r.Product).FirstOrCreate(&species, Species{Name: r.Product})

		// --- Insert Category ---
		var category Category
		db.Where("name = ?", r.Category).FirstOrCreate(&category, Category{Name: r.Category})

		// --- Insert Region ---
		var region Region
		db.Where("region = ?", r.Country).FirstOrCreate(&region, Region{Region: r.Country})

		// --- Insert Seafood ---
		seafood := Seafood{
			SpeciesID:  species.ID,
			RegionID:   region.ID,
			CategoryID: category.ID,
			PriceUnit:  priceUnit,
		}
		db.Create(&seafood)

		// --- Insert Price ---
		price := Price{
			SeafoodID: seafood.ID,
			Date:      parsedDate,
			Price:     priceValue,
		}
		db.Create(&price)
	}
	fmt.Println("✅ Seeding completed successfully!")
}

/// ---------- HELPER FUNCTIONS ---------- ///

func LoadCSV(filename string) ([]*PriceCSV, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []*PriceCSV
	if err := gocsv.UnmarshalFile(file, &records); err != nil {
		return nil, err
	}
	return records, nil
}
