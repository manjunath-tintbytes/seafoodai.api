package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ----------- Custom Date Type -----------
type CustomDate time.Time

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(cd)
	formatted := fmt.Sprintf("\"%s\"", t.Format("January 2, 2006")) // "October 16, 2025"
	return []byte(formatted), nil
}

// ----------- Response Struct -----------
type MarketSignalResponse struct {
	Title         string     `json:"title"`
	PublishedDate CustomDate `json:"published_date"`
}

// ----------- Handler -----------
func GetMarketSignals(db *gorm.DB) gin.HandlerFunc {
	stmt := `
		SELECT 
			title,
			published_date
		FROM market_signals
		WHERE deleted_at IS NULL
		ORDER BY published_date DESC, title ASC
		LIMIT $1
	`

	return func(c *gin.Context) {
		limit := 100
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}

		var results []MarketSignalResponse
		if err := db.Raw(stmt, limit).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}
