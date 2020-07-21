package amqpfx

import (
	"context"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/consul/api"
	"github.com/streadway/amqp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(
	provideAMQPConnection,
	provideAMQPChannel,
)

func provideAMQPConnection(lc fx.Lifecycle, sugar *zap.SugaredLogger, kv *api.KV) *amqp.Connection {
	conf, err := newConfig(kv)
	if err != nil {
		sugar.Fatal(err)
	}

	conn, err := amqp.Dial(conf.url)
	if err != nil {
		sugar.Fatalw("Failed to dial amqp", "url", conf.url, "error", err)
	}

	lc.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return conn
}

func provideAMQPChannel(lc fx.Lifecycle, sugar *zap.SugaredLogger, conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	if err != nil {
		sugar.Fatalw("Failed to new channel", "error", err)
	}

	lc.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(ctx context.Context) error {
			return channel.Close()
		},
	})

	return channel
}

type config struct {
	url string
}

func newConfig(kv *api.KV) (config, error) {
	pair, _, err := kv.Get("amqp", nil)
	if err != nil {
		return config{}, fmt.Errorf("consul kv failed to get key %s: %w", "amqp", err)
	}

	url, err := jsonparser.GetString(pair.Value, "url")
	if err != nil {
		return config{}, fmt.Errorf("failed to get get key %s: %w", "url", err)
	}

	return config{
		url: url,
	}, nil
}
