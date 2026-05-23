package main

import (
	"context"
	"fmt"
	"go-practice/greetings"
	"go-practice/internal/album"
	mysql_db "go-practice/internal/db"
	"go-practice/internal/httpserver"
	"go-practice/internal/storage/memory"
	"go-practice/internal/storage/mysql"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

var albums = []album.Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Calloway", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	fmt.Println("hello world")
	fmt.Println(greetings.Hello("Rohan"))
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		})
	defer rdb.Close()

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	limiter := redis_rate.NewLimiter(rdb)

	var limit atomic.Pointer[redis_rate.Limit]
	limit.Store(&redis_rate.Limit{
		Rate:   10,
		Burst:  20,
		Period: 1 * time.Minute,
	})
	httpserver.StartConfigWatcher(context.Background(), rdb, &limit, 10*time.Second)

	var store album.Store

	dsn := os.Getenv("MYSQL_DSN")
	if dsn != "" {
		fmt.Println("Using MySQL")
		db, err := mysql_db.InitDB(dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		store = mysql.NewStore(db)
		defer store.Close()
	} else {
		fmt.Println("Using Memory")
		store = memory.NewStore(albums)
		defer store.Close()
	}

	router := httpserver.NewRouter(store, limiter, &limit, rdb)

	fmt.Println("Server is running on port 8080")
	router.Run(":8080")
}
