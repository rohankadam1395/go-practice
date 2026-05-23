package httpserver

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func StartConfigWatcher(ctx context.Context, rdb *redis.Client, ptr *atomic.Pointer[redis_rate.Limit], interval time.Duration) {
	fmt.Println("starting config watcher")
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Println("checking rate limit config")
				vals, err := rdb.HGetAll(ctx, "rate_limit").Result()
				if err != nil {
					fmt.Printf("error getting rate limit from redis: %v", err)
					continue
				}
				if len(vals) == 0 {
					fmt.Println("no rate limit config found")
					continue
				}
				rate, err := strconv.Atoi(vals["rate"])
				if err != nil {
					fmt.Printf("error parsing rate: %v", err)
					continue
				}
				burst, err := strconv.Atoi(vals["burst"])
				if err != nil {
					fmt.Printf("error parsing burst: %v", err)
					continue
				}
				period, err := time.ParseDuration(vals["period"])
				if err != nil {
					fmt.Printf("error parsing period: %v", err)
					continue
				}
				limit := redis_rate.Limit{
					Rate:   rate,
					Burst:  burst,
					Period: period,
				}
				ptr.Store(&limit)
				fmt.Println("rate limit config updated")
			case <-ctx.Done():
				return
			}
		}
	}()
}
