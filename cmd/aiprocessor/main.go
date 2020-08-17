package main

import (
	"log"
	"os"

	"github.com/haunt98/colossus/internal/ai"
	"github.com/haunt98/colossus/internal/pkg/fx/aifx"
	"github.com/haunt98/colossus/internal/pkg/fx/amqpfx"
	"github.com/haunt98/colossus/internal/pkg/fx/bucketfx"
	"github.com/haunt98/colossus/internal/pkg/fx/cachefx"
	"github.com/haunt98/colossus/internal/pkg/fx/consulfx"
	"github.com/haunt98/colossus/internal/pkg/fx/miniofx"
	"github.com/haunt98/colossus/internal/pkg/fx/queuefx"
	"github.com/haunt98/colossus/internal/pkg/fx/redisfx"
	"github.com/haunt98/colossus/internal/pkg/fx/zapfx"
	"github.com/haunt98/colossus/pkg/bucket"
	"github.com/haunt98/colossus/pkg/cache"
	"github.com/haunt98/colossus/pkg/queue"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	project := os.Getenv("PROJECT")
	if project == "" {
		log.Fatal("Empty PROJECT")
	}

	fx.New(
		zapfx.Module,
		consulfx.Module,
		redisfx.Module,
		amqpfx.Module,
		miniofx.Module,
		fx.Provide(fx.Annotated{
			Name:   "project",
			Target: cachefx.InjectCache(project),
		}),
		fx.Provide(queuefx.InjectQueue(project)),
		fx.Provide(fx.Annotated{
			Name:   "storage",
			Target: cachefx.InjectCache("storage"),
		}),
		fx.Provide(bucketfx.InjectBucket("storage")),
		fx.Provide(aifx.InjectCMDConfig(project)),
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
	CMDConfig    ai.CMDConfig
}

func register(p params) {
	processor := ai.NewProcessor(p.Sugar, p.ProjectCache, p.Queue, p.StorageCache, p.Bucket, p.CMDConfig)
	if err := processor.Consume(); err != nil {
		p.Sugar.Error(err)
	}
}
