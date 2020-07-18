package cachefx

import (
	"colossus/pkg/cache"
	"fmt"

	"github.com/hashicorp/consul/api"

	"github.com/buger/jsonparser"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type ProvideCacheFn func(sugar *zap.SugaredLogger, client *redis.Client, kv *api.KV) *cache.Cache

func InjectCache(project string) ProvideCacheFn {
	return func(sugar *zap.SugaredLogger, client *redis.Client, kv *api.KV) *cache.Cache {
		conf, err := newConfig(kv, project)
		if err != nil {
			sugar.Fatal(err)
		}

		c := cache.NewCache(client, cache.PrefixKey(conf.prefix), cache.DefaultExpiration())

		return c
	}
}

type config struct {
	prefix string
}

func newConfig(kv *api.KV, project string) (config, error) {
	pair, _, err := kv.Get(project, nil)
	if err != nil {
		return config{}, fmt.Errorf("consul kv failed to get key %s: %w", project, err)
	}

	prefix, err := jsonparser.GetString(pair.Value, "cache", "prefix")
	if err != nil {
		return config{}, fmt.Errorf("failed to get key %s: %w", "cache.prefix", err)
	}

	return config{
		prefix: prefix,
	}, nil
}
