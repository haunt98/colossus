package miniofx

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/consul/api"
	"github.com/minio/minio-go/v6"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(
	provideMinio,
)

func provideMinio(sugar *zap.SugaredLogger, kv *api.KV) *minio.Client {
	conf, err := newConfig(kv)
	if err != nil {
		sugar.Fatal(err)
	}

	client, err := minio.New(conf.endpoint, conf.accessKeyID, conf.secretAccessKey, conf.useSSL)
	if err != nil {
		sugar.Fatalw("Failed to new minio client", "error", err)
	}

	return client
}

type config struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
}

func newConfig(kv *api.KV) (config, error) {
	pair, _, err := kv.Get("minio", nil)
	if err != nil {
		return config{}, fmt.Errorf("consul kv failed to get key %s: %w", "minio", err)
	}

	endpoint, err := jsonparser.GetString(pair.Value, "endpoint")
	if err != nil {
		return config{}, fmt.Errorf("failed to get key %s: %w", "endpoint", err)
	}

	accessKeyID, err := jsonparser.GetString(pair.Value, "access_key_id")
	if err != nil {
		return config{}, fmt.Errorf("failed to get key %s: %w", "access_key_id", err)
	}

	secretAccessKey, err := jsonparser.GetString(pair.Value, "secret_access_key")
	if err != nil {
		return config{}, fmt.Errorf("failed to get key %s: %w", "secret_access_key", err)
	}

	useSSL, err := jsonparser.GetBoolean(pair.Value, "use_ssl")
	if err != nil {
		return config{}, fmt.Errorf("failed to get key %s: %w", "use_ssl", err)
	}

	return config{
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		useSSL:          useSSL,
	}, nil
}
