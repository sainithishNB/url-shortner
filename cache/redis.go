package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetURL(rdb *redis.Client, code string) (string, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, code).Result()
	return val, err
}
func SetURL(rdb *redis.Client, code string, longURL string, ttl time.Duration) {
	ctx := context.Background()
	rdb.Set(ctx, code, longURL, ttl)
}
func DelURL(rdb *redis.Client, code string) {
	ctx := context.Background()
	rdb.Del(ctx, code)
}
