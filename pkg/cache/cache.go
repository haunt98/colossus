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

func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	key = c.generateKeyFn(key)

	data, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("client failed to get: %w", err)
	}

	return data, nil
}

func (c *Cache) SetString(ctx context.Context, key string, value string) error {
	key = c.generateKeyFn(key)

	if err := c.redisClient.Set(ctx, key, value, c.expiration).Err(); err != nil {
		return fmt.Errorf("client failed to set: %w", err)
	}

	return nil
}

func (c *Cache) GetJSON(ctx context.Context, key string, value interface{}) error {
	key = c.generateKeyFn(key)

	data, err := c.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return fmt.Errorf("redis client failed to get: %w", err)
	}

	if easyValue, ok := value.(easyjson.Unmarshaler); ok {
		if err := easyjson.Unmarshal(data, easyValue); err != nil {
			return fmt.Errorf("json failed to unmarshal: %w", err)
		}

		return nil
	} else {
		if err := json.Unmarshal(data, value); err != nil {
			return fmt.Errorf("json failed to unmarshal: %w", err)
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
			return fmt.Errorf("json failed to marshal: %w", err)
		}
	} else {
		if data, err = json.Marshal(value); err != nil {
			return fmt.Errorf("json failed to marshal: %w", err)
		}
	}

	if err := c.redisClient.Set(ctx, key, data, c.expiration).Err(); err != nil {
		return fmt.Errorf("redis client failed to set: %w", err)
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
