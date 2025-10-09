package handlers

import (
	"fmt"
	"net/http"
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

// ------------------ Handlers ------------------

func GetMarketPricesOptimized(db *gorm.DB) gin.HandlerFunc {
	stmt := `
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
			WHERE p.deleted_at IS NULL 
			  AND s.deleted_at IS NULL 
			  AND sp.deleted_at IS NULL 
			  AND r.deleted_at IS NULL
			ORDER BY s.species_id, s.region_id, p.date DESC
		),
		latest_limited AS (
			SELECT * FROM (
				SELECT *, ROW_NUMBER() OVER (ORDER BY date DESC) as rn 
				FROM latest_per_species_region
			) ranked WHERE rn <= $1
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
		ORDER BY ll.date DESC`

	return func(c *gin.Context) {
		limit := 50
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}

		var results []MarketPriceResult
		err := db.Raw(stmt, limit).Scan(&results).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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

		c.JSON(http.StatusOK, marketPrices)
	}
}
