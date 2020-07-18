package consulfx

import (
	"os"

	"github.com/hashicorp/consul/api"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(
	provideClient,
	provideAgent,
	provideKV,
)

const consulAdressEnv = "CONSUL_ADDRESS"

func provideClient(sugar *zap.SugaredLogger) *api.Client {
	address := os.Getenv(consulAdressEnv)
	if address == "" {
		sugar.Fatalf("Empty %s", consulAdressEnv)
	}

	client, err := api.NewClient(&api.Config{
		Address: address,
	})
	if err != nil {
		sugar.Fatalw("Failed to new consul client", "error", err)
	}

	return client
}

func provideAgent(client *api.Client) *api.Agent {
	return client.Agent()
}

func provideKV(client *api.Client) *api.KV {
	return client.KV()
}
