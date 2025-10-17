package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/handlers"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/market-prices", handlers.GetMarketPricesOptimized(db))
	r.GET("/landings", handlers.GetLandings(db))
	r.GET("/market-signals", handlers.GetMarketSignals(db))
}
