package bucketfx

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/consul/api"
	"github.com/haunt98/colossus/pkg/bucket"
	"github.com/minio/minio-go/v6"
	"go.uber.org/zap"
)

type ProvideBucketFn func(sugar *zap.SugaredLogger, client *minio.Client, kv *api.KV) *bucket.Bucket

func InjectBucket(project string) ProvideBucketFn {
	return func(sugar *zap.SugaredLogger, client *minio.Client, kv *api.KV) *bucket.Bucket {
		conf, err := newConfig(kv, project)
		if err != nil {
			sugar.Error(err)
		}

		b, err := bucket.NewBucket(client, conf.name, bucket.DefaultPresignedExpiration())
		if err != nil {
			sugar.Fatal(err)
		}

		return b
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

	name, err := jsonparser.GetString(pair.Value, "bucket", "name")
	if err != nil {
		return config{}, fmt.Errorf("failed to get key %s: %w", "bucket.name", err)
	}

	return config{
		name: name,
	}, nil
}
