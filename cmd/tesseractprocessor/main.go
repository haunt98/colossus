package main

import (
	"colossus/internal/ai"
	"colossus/internal/pkg/fx/bucketfx"
	"colossus/internal/pkg/fx/cachefx"
	"colossus/internal/pkg/fx/consulfx"
	"colossus/internal/pkg/fx/miniofx"
	"colossus/internal/pkg/fx/queuefx"
	"colossus/internal/pkg/fx/rabbitmqfx"
	"colossus/internal/pkg/fx/redisfx"
	"colossus/internal/pkg/fx/zapfx"
	"colossus/pkg/bucket"
	"colossus/pkg/cache"
	"colossus/pkg/queue"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	project := "tesseract"

	fx.New(
		zapfx.Module,
		consulfx.Module,
		redisfx.Module,
		rabbitmqfx.Module,
		miniofx.Module,
		fx.Provide(fx.Annotated{
			Name:   project,
			Target: cachefx.InjectCache(project),
		}),
		fx.Provide(queuefx.InjectQueue(project)),
		fx.Provide(fx.Annotated{
			Name:   "storage",
			Target: cachefx.InjectCache("storage"),
		}),
		fx.Provide(bucketfx.InjectBucket("storage")),
		fx.Invoke(register),
	).Run()
}

type params struct {
	fx.In

	Sugar        *zap.SugaredLogger
	ProjectCache *cache.Cache `name:"project"`
	Queue        *queue.Queue
	StorageCache *cache.Cache `name:"storage"`
	Bucket       *bucket.Bucket
}

func register(p params) {
	processor := ai.NewProcessor(p.Sugar, p.ProjectCache, p.Queue, p.StorageCache, p.Bucket, ai.CMDConfig{})
	if err := processor.Consume(); err != nil {
		p.Sugar.Error(err)
	}
}
