package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"order-service/config"
)

type Client struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedis(cfg config.Config) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost + ":" + cfg.RedisPort,
	})
	return &Client{rdb: rdb, ctx: context.Background()}
}

func (c *Client) Get(key string, dest interface{}) error {
	val, err := c.rdb.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *Client) Set(key string, value interface{}, ttl int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	duration := time.Duration(ttl) * time.Second
	return c.rdb.Set(c.ctx, key, data, duration).Err()
}

func (c *Client) Del(key string) error {
	return c.rdb.Del(c.ctx, key).Err()
}
