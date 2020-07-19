package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mailru/easyjson"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	redisClient   *redis.Client
	generateKeyFn GenerateKeyFn
	expiration    time.Duration
}

type GenerateKeyFn func(key string) string

func NewCache(redisClient *redis.Client, generateKeyFn GenerateKeyFn, expiration time.Duration) *Cache {
	return &Cache{
		redisClient:   redisClient,
		generateKeyFn: generateKeyFn,
		expiration:    expiration,
	}
}

func (c *Cache) GetJSON(ctx context.Context, key string, value interface{}) error {
	key = c.generateKeyFn(key)

	data, err := c.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return fmt.Errorf("redis client failed to get key %s: %w", key, err)
	}

	if easyValue, ok := value.(easyjson.Unmarshaler); ok {
		if err := easyjson.Unmarshal(data, easyValue); err != nil {
			return fmt.Errorf("json failed to unmarshal %s: %w", data, err)
		}

		return nil
	} else {
		if err := json.Unmarshal(data, value); err != nil {
			return fmt.Errorf("json failed to unmarshal %s: %w", data, err)
		}

		return nil
	}
}

func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}) error {
	key = c.generateKeyFn(key)

	var data []byte
	var err error
	if easyValue, ok := value.(easyjson.Marshaler); ok {
		if data, err = easyjson.Marshal(easyValue); err != nil {
			return fmt.Errorf("json failed to marshal %v: %w", value, err)
		}
	} else {
		if data, err = json.Marshal(value); err != nil {
			return fmt.Errorf("json failed to marshal %v: %w", value, err)
		}
	}

	if err := c.redisClient.Set(ctx, key, data, c.expiration).Err(); err != nil {
		return fmt.Errorf("redis client failed to set key %s: %w", key, err)
	}

	return nil
}

func PrefixKey(prefix string) GenerateKeyFn {
	return func(key string) string {
		return prefix + ":" + key
	}
}

func DefaultExpiration() time.Duration {
	return time.Hour * 12
}
