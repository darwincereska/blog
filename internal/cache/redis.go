package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

// CreateCache creates a new Redis Cache
func CreateCache(host string, port int, db int) *Cache {
	address := fmt.Sprintf("%s:%d", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       db,
	})

	return &Cache{
		Client: rdb,
	}
}

// Set sets a key value in the cache with an expiry
func (c *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

// Get returns a key value from the cache if there is any
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

// Clear clears specific keys from the cache
func (c *Cache) Clear(ctx context.Context, keys ...string) error {
	return c.Client.Del(ctx, keys...).Err()
}

// ClearPattern clears keys by pattern (e.g. "user:*")
func (c *Cache) ClearPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := c.Client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			// Unlink instead of Del for large datasets because it's non-blocking
			if err := c.Client.Unlink(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

// Flush clears entire current db
func (c *Cache) Flush(ctx context.Context) error {
	return c.Client.FlushDB(ctx).Err()
}
