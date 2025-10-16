package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/manjunath-tintbytes/seafoodai.api/internal/config"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/database"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnv()
	db := database.SetupDB()

	filePath := "./cmd/seeder/landings/FOSS_landings.xlsx"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("❌ Failed to open Excel file: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("❌ Failed to read sheet: %v", err)
	}

	if len(rows) < 2 {
		log.Fatal("❌ Excel file has no data rows.")
	}

	// Cache maps to prevent repeated DB lookups
	nameCache := make(map[string]uint)
	portCache := make(map[string]uint)

	tx := db.Begin() // start transaction

	for i, row := range rows[1:] {
		if len(row) < 10 {
			log.Printf("⚠️ Row %d is too short, skipping", i+2)
			continue
		}

		year := parseInt(row[0])
		regionName := row[1]
		nmfsName := row[2]
		pounds := parseFloat(row[3])
		dollars := parseFloat(row[4])
		scientificName := row[8]
		metricTons := parseFloat(row[9])

		// Skip incomplete or zero-value rows
		if pounds <= 0 || dollars <= 0 || metricTons <= 0 {
			continue
		}

		if nmfsName == "" || regionName == "" {
			continue
		}

		// --- LANDING NAMES ---
		nameKey := nmfsName + "|" + scientificName
		var landingNameID uint
		if id, ok := nameCache[nameKey]; ok {
			landingNameID = id
		} else {
			var existing models.LandingName
			if err := tx.Where("nmfs_name = ? AND scientific_name = ?", nmfsName, scientificName).
				First(&existing).Error; err == gorm.ErrRecordNotFound {
				newName := models.LandingName{
					NMFSName:       nmfsName,
					ScientificName: scientificName,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				if err := tx.Create(&newName).Error; err != nil {
					log.Printf("❌ Failed to insert landing name (row %d): %v", i+2, err)
					continue
				}
				landingNameID = newName.ID
			} else {
				landingNameID = existing.ID
			}
			nameCache[nameKey] = landingNameID
		}

		// --- LANDING PORTS ---
		portKey := regionName
		var landingPortID uint
		if id, ok := portCache[portKey]; ok {
			landingPortID = id
		} else {
			var existing models.LandingPort
			if err := tx.Where("region_name = ?", regionName).
				First(&existing).Error; err == gorm.ErrRecordNotFound {
				newPort := models.LandingPort{
					RegionName: regionName,
					Port:       "", // fill later if needed
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				if err := tx.Create(&newPort).Error; err != nil {
					log.Printf("❌ Failed to insert landing port (row %d): %v", i+2, err)
					continue
				}
				landingPortID = newPort.ID
			} else {
				landingPortID = existing.ID
			}
			portCache[portKey] = landingPortID
		}

		// --- LANDINGS ---
		landing := models.Landing{
			Year:          year,
			LandingPortID: landingPortID,
			LandingNameID: landingNameID,
			Pounds:        pounds,
			Dollars:       dollars,
			MetricTons:    metricTons,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := tx.Create(&landing).Error; err != nil {
			log.Printf("❌ Failed to insert landing (row %d): %v", i+2, err)
			continue
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("❌ Transaction commit failed: %v", err)
	}

	fmt.Println("✅ Landings data imported successfully without duplicates!")
}

func parseFloat(s string) float64 {
	var v float64
	fmt.Sscanf(s, "%f", &v)
	return v
}

func parseInt(s string) int {
	var v int
	fmt.Sscanf(s, "%d", &v)
	return v
}
