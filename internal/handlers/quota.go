package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ----------- Response Struct -----------
type QuotaResponse struct {
	Date           string `json:"date"`
	ProductName    string `json:"product_name"`
	RemainingQuota string `json:"remaining_quota"`
}

// ----------- Handler -----------
func GetQuotas(db *gorm.DB) gin.HandlerFunc {
	stmt := `
		SELECT 
			TO_CHAR(date, 'YYYY-MM-DD') AS date,
			product_name,
			CONCAT(remaining_quota, '%') AS remaining_quota
		FROM quota
		WHERE deleted_at IS NULL
		  AND date = (SELECT MAX(date) FROM quota WHERE deleted_at IS NULL)
		ORDER BY product_name ASC
		LIMIT $1
	`

	return func(c *gin.Context) {
		limit := 100
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}

		var results []QuotaResponse
		if err := db.Raw(stmt, limit).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}
