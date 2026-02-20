package cache

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/logging"
)

type Cache struct {
	Client  *redis.Client
	Enabled atomic.Bool
	logger *log.Logger
}

// CreateCache creates a new Redis Cache
func CreateCache(host string, port int, db int) *Cache {
	address := fmt.Sprintf("%s:%d", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:       address,
		Password:   "",
		DB:         db,
		MaxRetries: -1,
	})

	c := &Cache{
		Client: rdb,
		logger: log.NewWithOptions(os.Stderr, log.Options{
			ReportTimestamp: true,	
			ReportCaller: false,
			Prefix: "REDIS",
		}),
	}

	logging.Disable() // Disable redis internal logging

	// Initial check to set the state
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		c.logger.Warn("Redis is offline. Starting in disabled mode.", "error", err)
		c.Enabled.Store(false)
	} else {
		c.Enabled.Store(true)
	}

	// Start background health checker
	go c.startHealthCheck()

	return c
}

// startHealthCheck periodically checks the status of redis and changes status
func (c *Cache) startHealthCheck() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := c.Client.Ping(ctx).Err()
		cancel()

		if err != nil {
			if c.Enabled.Load() {
				c.logger.Error("Redis connection lost. Disabling cache.")
				c.Enabled.Store(false)
			}
		} else {
			if !c.Enabled.Load() {
				c.logger.Info("Redis connection restored. Enabling cache.")
				c.Enabled.Store(true)
			}
		}
	}
}

// Set sets a key value in the cache with an expiry
func (c *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if !c.Enabled.Load() {
		return nil
	}

	return c.Client.Set(ctx, key, value, expiration).Err()
}

// Get returns a key value from the cache if there is any
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	if !c.Enabled.Load() {
		return "", nil
	}

	val, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		// Handle key not found
		if err == redis.Nil {
			// Return empty key and return descriptive error
			return "", fmt.Errorf("key: %s not found. Error: %w", key, err)
		}
		// Handle any other errors
		return "", fmt.Errorf("redis GET %s failed: %w", key, err)
	}
	return val, nil
}

// Clear clears specific keys from the cache
func (c *Cache) Clear(ctx context.Context, keys ...string) error {
	if !c.Enabled.Load() {
		return nil
	}

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
	if !c.Enabled.Load() {
		return fmt.Errorf("redis is in disabled mode")
	}

	return c.Client.FlushDB(ctx).Err()
}
