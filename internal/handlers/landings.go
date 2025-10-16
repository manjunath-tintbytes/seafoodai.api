package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ----------- Response Struct -----------
type LandingResponse struct {
	Year       int     `json:"year"`
	RegionName string  `json:"region"`
	NMFSName   string  `json:"nmfs_name"`
	Pounds     float64 `json:"pounds"`
	Dollars    float64 `json:"dollars"`
	MetricTons float64 `json:"metric_tons"`
}

// ----------- Handler -----------
func GetLandings(db *gorm.DB) gin.HandlerFunc {
	stmt := `
		SELECT 
			l.year,
			lp.region_name,
			ln.nmfs_name,
			COALESCE(l.pounds, 0) AS pounds,
			COALESCE(l.dollars, 0) AS dollars,
			COALESCE(l.metric_tons, 0) AS metric_tons
		FROM landings l
		JOIN landing_ports lp ON l.landing_port_id = lp.id
		JOIN landing_names ln ON l.landing_name_id = ln.id
		WHERE l.deleted_at IS NULL 
		  AND lp.deleted_at IS NULL 
		  AND ln.deleted_at IS NULL
		ORDER BY l.year DESC, lp.region_name ASC, ln.nmfs_name ASC
		LIMIT $1
	`

	return func(c *gin.Context) {
		limit := 100
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}

		var results []LandingResponse
		if err := db.Raw(stmt, limit).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}
