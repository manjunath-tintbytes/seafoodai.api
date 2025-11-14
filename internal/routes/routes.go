package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/handlers"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/", handlers.WelcomeHandler(db))
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.POST("/forgot-password", handlers.ForgotPassword)
	r.POST("/reset-password", handlers.ResetPassword)

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
