package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/utils"
	"gorm.io/gorm"
)

type MarketPrice struct {
	SpeciesSKU  string   `json:"species_sku"`
	Origin      string   `json:"origin"`
	Price       float64  `json:"price"`
	PriceUnit   string   `json:"price_unit"`
	WeeklyTrend *float64 `json:"weekly_trend"`
	YoY         *float64 `json:"yoy"`
}

type MarketPriceResult struct {
	ID           uint      `json:"id"`
	Price        float64   `json:"price"`
	Date         time.Time `json:"date"`
	PriceUnit    string    `json:"price_unit"`
	SpeciesName  string    `json:"species_name"`
	RegionName   string    `json:"region_name"`
	SpeciesID    uint      `json:"species_id"`
	RegionID     uint      `json:"region_id"`
	WeekAgoPrice *float64  `json:"week_ago_price"`
	YearAgoPrice *float64  `json:"year_ago_price"`
}

type PaginatedResponse struct {
	Data       []MarketPrice `json:"data"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalCount int64         `json:"total_count"`
	TotalPages int           `json:"total_pages"`
}

// ------------------ Handlers ------------------

func GetMarketPricesOptimized(db *gorm.DB) gin.HandlerFunc {
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
		speciesName := strings.TrimSpace(c.Query("species"))
		regionName := strings.TrimSpace(c.Query("region"))

		// Build WHERE clause for filters
		var filterConditions []string
		var filterArgs []interface{}
		argIndex := 1

		baseWhereClause := `p.deleted_at IS NULL
			  AND s.deleted_at IS NULL
			  AND sp.deleted_at IS NULL
			  AND r.deleted_at IS NULL`

		if speciesName != "" {
			filterConditions = append(filterConditions, fmt.Sprintf("sp.name ILIKE $%d", argIndex))
			filterArgs = append(filterArgs, "%"+speciesName+"%")
			argIndex++
		}

		if regionName != "" {
			filterConditions = append(filterConditions, fmt.Sprintf("r.region ILIKE $%d", argIndex))
			filterArgs = append(filterArgs, "%"+regionName+"%")
			argIndex++
		}

		whereClause := baseWhereClause
		if len(filterConditions) > 0 {
			whereClause += " AND " + strings.Join(filterConditions, " AND ")
		}

		// Count query with filters
		countStmt := fmt.Sprintf(`
			SELECT COUNT(DISTINCT (s.species_id, s.region_id))
			FROM prices p
			JOIN seafoods s ON p.seafood_id = s.id
			JOIN species sp ON s.species_id = sp.id
			JOIN regions r ON s.region_id = r.id
			WHERE %s`, whereClause)

		// Main query with filters
		stmt := fmt.Sprintf(`
			WITH latest_per_species_region AS (
				SELECT DISTINCT ON (s.species_id, s.region_id)
					p.id,
					p.price,
					p.date,
					s.price_unit,
					sp.name as species_name,
					r.region as region_name,
					s.species_id,
					s.region_id
				FROM prices p
				JOIN seafoods s ON p.seafood_id = s.id
				JOIN species sp ON s.species_id = sp.id
				JOIN regions r ON s.region_id = r.id
				WHERE %s
				ORDER BY s.species_id, s.region_id, p.date DESC
			),
			latest_limited AS (
				SELECT * FROM latest_per_species_region
				ORDER BY date DESC
				LIMIT $%d OFFSET $%d
			),
			week_ago_prices AS (
				SELECT DISTINCT ON (s.species_id, s.region_id)
					s.species_id,
					s.region_id,
					p.price as week_ago_price
				FROM latest_limited ll
				JOIN seafoods s ON s.species_id = ll.species_id AND s.region_id = ll.region_id
				JOIN prices p ON p.seafood_id = s.id
				WHERE p.date < ll.date - INTERVAL '6 days'
				  AND p.deleted_at IS NULL
				  AND s.deleted_at IS NULL
				ORDER BY s.species_id, s.region_id, p.date DESC
			),
			year_ago_prices AS (
				SELECT DISTINCT ON (s.species_id, s.region_id)
					s.species_id,
					s.region_id,
					p.price as year_ago_price
				FROM latest_limited ll
				JOIN seafoods s ON s.species_id = ll.species_id AND s.region_id = ll.region_id
				JOIN prices p ON p.seafood_id = s.id
				WHERE p.date >= make_date(EXTRACT(YEAR FROM ll.date)::int - 1, 1, 1)
				  AND p.date < make_date(EXTRACT(YEAR FROM ll.date)::int, 1, 1)
				  AND p.deleted_at IS NULL
				  AND s.deleted_at IS NULL
				ORDER BY s.species_id, s.region_id, p.date ASC
			)
			SELECT
				ll.*,
				wap.week_ago_price,
				yap.year_ago_price
			FROM latest_limited ll
			LEFT JOIN week_ago_prices wap ON ll.species_id = wap.species_id AND ll.region_id = wap.region_id
			LEFT JOIN year_ago_prices yap ON ll.species_id = yap.species_id AND ll.region_id = yap.region_id
			ORDER BY ll.date DESC`, whereClause, argIndex, argIndex+1)

		// Calculate offset
		offset := (page - 1) * pageSize

		// Prepare arguments for queries
		queryArgs := append(filterArgs, pageSize, offset)

		// Get total count
		var totalCount int64
		err := db.Raw(countStmt, filterArgs...).Scan(&totalCount).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get paginated results
		var results []MarketPriceResult
		err = db.Raw(stmt, queryArgs...).Scan(&results).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Transform results
		marketPrices := make([]MarketPrice, len(results))
		for i, r := range results {
			marketPrices[i] = MarketPrice{
				SpeciesSKU:  r.SpeciesName,
				Origin:      r.RegionName,
				Price:       r.Price,
				PriceUnit:   r.PriceUnit,
				WeeklyTrend: utils.CalculateChange(r.Price, r.WeekAgoPrice),
				YoY:         utils.CalculateChange(r.Price, r.YearAgoPrice),
			}
		}

		// Calculate total pages
		totalPages := int(totalCount) / pageSize
		if int(totalCount)%pageSize != 0 {
			totalPages++
		}

		// Return paginated response
		c.JSON(http.StatusOK, PaginatedResponse{
			Data:       marketPrices,
			Page:       page,
			PageSize:   pageSize,
			TotalCount: totalCount,
			TotalPages: totalPages,
		})
	}
}
