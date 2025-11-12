package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/handlers"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/", handlers.WelcomeHandler(db))
	// r.GET("/market-prices", handlers.GetMarketPricesOptimized(db))
	// r.GET("/landings", handlers.GetLandings(db))
	// r.GET("/market-signals", handlers.GetMarketSignals(db))
	// r.GET("/quotas", handlers.GetQuotas(db))
	// Public routes
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)

	// Protected routes
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.GET("/market-prices", handlers.GetMarketPricesOptimized(db))
		protected.GET("/landings", handlers.GetLandings(db))
		protected.GET("/market-signals", handlers.GetMarketSignals(db))
		protected.GET("/quotas", handlers.GetQuotas(db))
	}
}
