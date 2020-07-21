package queuefx

import (
	"colossus/pkg/queue"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/consul/api"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type ProvideQueueFn func(sugar *zap.SugaredLogger, channel *amqp.Channel, kv *api.KV) *queue.Queue

func InjectQueue(project string) ProvideQueueFn {
	return func(sugar *zap.SugaredLogger, channel *amqp.Channel, kv *api.KV) *queue.Queue {
		conf, err := newConfig(kv, project)
		if err != nil {
			sugar.Fatal(err)
		}

		q, err := queue.NewQueue(channel, conf.name)
		if err != nil {
			sugar.Fatal(err)
		}

		return q
	}
}

type config struct {
	name string
}

func newConfig(kv *api.KV, project string) (config, error) {
	pair, _, err := kv.Get(project, nil)
	if err != nil {
		return config{}, fmt.Errorf("consul kv failed to get key %s: %w", project, err)
	}

	name, err := jsonparser.GetString(pair.Value, "queue", "name")
	if err != nil {
		return config{}, fmt.Errorf("failed to get key %s: %w", "queue.name", err)
	}

	return config{
		name: name,
	}, nil
}
