package handlers

import (
	"fmt"
	"net/http"
	"strings"

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

type LandingsPaginatedResponse struct {
	Data       []LandingResponse `json:"data"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalCount int64             `json:"total_count"`
	TotalPages int               `json:"total_pages"`
}

// ----------- Handler -----------

func GetLandings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse pagination parameters
		page := 1
		pageSize := 20

		if p := c.Query("page"); p != "" {
			fmt.Sscanf(p, "%d", &page)
			if page < 1 {
				page = 1
			}
		}

		if ps := c.Query("page_size"); ps != "" {
			fmt.Sscanf(ps, "%d", &pageSize)
			if pageSize < 1 {
				pageSize = 20
			}
		}

		// Parse filter parameters
		yearStr := strings.TrimSpace(c.Query("year"))
		regionName := strings.TrimSpace(c.Query("region"))
		nmfsName := strings.TrimSpace(c.Query("name"))

		// Build WHERE clause for filters
		var filterConditions []string
		var filterArgs []interface{}
		argIndex := 1

		baseWhereClause := `l.deleted_at IS NULL 
			  AND lp.deleted_at IS NULL 
			  AND ln.deleted_at IS NULL`

		if yearStr != "" {
			var year int
			if _, err := fmt.Sscanf(yearStr, "%d", &year); err == nil {
				filterConditions = append(filterConditions, fmt.Sprintf("l.year = $%d", argIndex))
				filterArgs = append(filterArgs, year)
				argIndex++
			}
		}

		if regionName != "" {
			filterConditions = append(filterConditions, fmt.Sprintf("lp.region_name ILIKE $%d", argIndex))
			filterArgs = append(filterArgs, "%"+regionName+"%")
			argIndex++
		}

		if nmfsName != "" {
			filterConditions = append(filterConditions, fmt.Sprintf("ln.nmfs_name ILIKE $%d", argIndex))
			filterArgs = append(filterArgs, "%"+nmfsName+"%")
			argIndex++
		}

		whereClause := baseWhereClause
		if len(filterConditions) > 0 {
			whereClause += " AND " + strings.Join(filterConditions, " AND ")
		}

		// Count query with filters
		countStmt := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM landings l
			JOIN landing_ports lp ON l.landing_port_id = lp.id
			JOIN landing_names ln ON l.landing_name_id = ln.id
			WHERE %s`, whereClause)

		// Main query with filters and pagination
		stmt := fmt.Sprintf(`
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
			WHERE %s
			ORDER BY l.year DESC, lp.region_name ASC, ln.nmfs_name ASC
			LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

		// Calculate offset
		offset := (page - 1) * pageSize

		// Prepare arguments for queries
		queryArgs := append(filterArgs, pageSize, offset)

		// Get total count
		var totalCount int64
		if err := db.Raw(countStmt, filterArgs...).Scan(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get paginated results
		var results []LandingResponse
		if err := db.Raw(stmt, queryArgs...).Scan(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Calculate total pages
		totalPages := int(totalCount) / pageSize
		if int(totalCount)%pageSize != 0 {
			totalPages++
		}

		// Return paginated response
		c.JSON(http.StatusOK, LandingsPaginatedResponse{
			Data:       results,
			Page:       page,
			PageSize:   pageSize,
			TotalCount: totalCount,
			TotalPages: totalPages,
		})
	}
}
