package httpserver

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func RateLimit(limiter *redis_rate.Limiter, limit *atomic.Pointer[redis_rate.Limit]) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		currentLimit := limit.Load()
		// fmt.Println("rate limit", currentLimit)
		res, err := limiter.Allow(c.Request.Context(), key, *currentLimit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if res.Allowed == 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
