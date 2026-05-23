package httpserver

import (
	"go-practice/internal/album"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func NewRouter(store album.Store, limiter *redis_rate.Limiter, limit *atomic.Pointer[redis_rate.Limit], rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(RateLimit(limiter, limit))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World"})
	})

	router.GET("/albums", func(c *gin.Context) {
		getAlbums(c, store)
	})

	router.POST("/albums", func(c *gin.Context) {
		postAlbums(c, store)
	})

	router.GET("/album/:id", func(c *gin.Context) {
		getAlbumByID(c, store)
	})

	router.PUT("/album/:id", func(c *gin.Context) {
		updateAlbum(c, store)
	})

	router.DELETE("/album/:id", func(c *gin.Context) {
		deleteAlbum(c, store)
	})

	router.POST("/admin/ratelimit", func(c *gin.Context) {
		updateRateLimit(c, rdb)
	})

	return router
}
