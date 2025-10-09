package main

import (
	"github.com/manjunath-tintbytes/seafoodai.api/internal/config"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/database"

	"log"
	"os"
)

func main() {
	config.LoadEnv()
	db := database.SetupDB()

	// Run migrations
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			if err := database.RunMigrations(db); err != nil {
				log.Fatal("Migration failed:", err)
			}
		case "rollback":
			if err := database.RollbackMigration(db); err != nil {
				log.Fatal("Rollback failed:", err)
			}
		case "rollback-to":
			if len(os.Args) < 3 {
				log.Fatal("Please provide migration ID")
			}
			if err := database.RollbackToMigration(db, os.Args[2]); err != nil {
				log.Fatal("Rollback failed:", err)
			}
		default:
			log.Println("Unknown command. Use: migrate, rollback, or rollback-to <migration_id>")
		}
	} else {
		// Your application logic here
		log.Println("Application started")
	}
}
