package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"sea-cucumber-trace/backend/internal/config"
)

type Client struct {
	rdb    *redis.Client
	prefix string
	ttl    time.Duration
}

func New(cfg *config.Config) (*Client, error) {
	if !cfg.RedisEnabled {
		return nil, nil
	}
	c := &Client{
		rdb: redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		}),
		prefix: cfg.RedisKeyPrefix,
		ttl:    time.Duration(cfg.CacheTTL) * time.Second,
	}
	if c.ttl <= 0 {
		c.ttl = 5 * time.Minute
	}
	if err := c.rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Enabled() bool {
	return c != nil && c.rdb != nil
}

func (c *Client) Close() error {
	if !c.Enabled() {
		return nil
	}
	return c.rdb.Close()
}

func (c *Client) GetJSON(ctx context.Context, key string, dest any) (bool, error) {
	if !c.Enabled() {
		return false, nil
	}
	raw, err := c.rdb.Get(ctx, c.key(key)).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(raw, dest); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) SetJSON(ctx context.Context, key string, value any) error {
	if !c.Enabled() {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, c.key(key), raw, c.ttl).Err()
}

func (c *Client) Delete(ctx context.Context, keys ...string) error {
	if !c.Enabled() || len(keys) == 0 {
		return nil
	}
	full := make([]string, 0, len(keys))
	for _, key := range keys {
		full = append(full, c.key(key))
	}
	return c.rdb.Del(ctx, full...).Err()
}

func (c *Client) FlushPrefix(ctx context.Context) error {
	if !c.Enabled() {
		return nil
	}

	var cursor uint64
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, c.prefix+"*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			return nil
		}
	}
}

func (c *Client) Ping(ctx context.Context) error {
	if !c.Enabled() {
		return nil
	}
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) key(key string) string {
	return c.prefix + key
}
