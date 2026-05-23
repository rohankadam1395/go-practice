package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type RateLimitRequest struct {
	Rate   int    `json:"rate"`
	Burst  int    `json:"burst"`
	Period string `json:"period"`
}

func updateRateLimit(c *gin.Context, rdb *redis.Client) {
	var dto RateLimitRequest
	if err := c.BindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	period, err := time.ParseDuration(dto.Period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := redis_rate.Limit{
		Rate:   dto.Rate,
		Burst:  dto.Burst,
		Period: period,
	}
	if err := rdb.HSet(c.Request.Context(), "rate_limit", "rate", limit.Rate, "burst", limit.Burst, "period", limit.Period.String()).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rate limit updated successfully", "limit": limit})
}
