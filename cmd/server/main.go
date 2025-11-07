package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/config"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/database"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/routes"
)

func main() {
	// Load environment variables
	config.LoadEnv()
	gin.SetMode(gin.ReleaseMode)

	// Connect to the database
	db := database.SetupDB()

	// Setup Gin router
	r := gin.Default()

	// Allow all origins
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Register routes
	routes.RegisterRoutes(r, db)

	// Start the server
	r.Run(":8080")
}
