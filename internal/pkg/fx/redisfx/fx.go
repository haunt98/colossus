package redisfx

import (
	"context"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(
	provideRedis,
)

func provideRedis(sugar *zap.SugaredLogger, kv *api.KV) *redis.Client {
	conf, err := newConfig(kv)
	if err != nil {
		sugar.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     conf.addr,
		Password: conf.password,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		sugar.Fatalw("Redis client failed to ping", "error", err)
	}

	return client
}

type config struct {
	addr     string
	password string
}

func newConfig(kv *api.KV) (config, error) {
	pair, _, err := kv.Get("redis", nil)
	if err != nil {
		return config{}, fmt.Errorf("consul kv failed to get key %s: %w", "redis", err)
	}

	addr, err := jsonparser.GetString(pair.Value, "addr")
	if err != nil {
		return config{}, fmt.Errorf("failed to get get key %s: %w", "addr", err)
	}

	password, err := jsonparser.GetString(pair.Value, "password")
	if err != nil {
		return config{}, fmt.Errorf("failed to get get key %s: %w", "password", err)
	}

	return config{
		addr:     addr,
		password: password,
	}, nil
}
