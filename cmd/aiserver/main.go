package main

import (
	"log"
	"os"

	"github.com/haunt98/colossus/internal/ai"
	"github.com/haunt98/colossus/internal/pkg/fx/amqpfx"
	"github.com/haunt98/colossus/internal/pkg/fx/bucketfx"
	"github.com/haunt98/colossus/internal/pkg/fx/cachefx"
	"github.com/haunt98/colossus/internal/pkg/fx/consulfx"
	"github.com/haunt98/colossus/internal/pkg/fx/grpcfx"
	"github.com/haunt98/colossus/internal/pkg/fx/miniofx"
	"github.com/haunt98/colossus/internal/pkg/fx/queuefx"
	"github.com/haunt98/colossus/internal/pkg/fx/redisfx"
	"github.com/haunt98/colossus/internal/pkg/fx/zapfx"
	"github.com/haunt98/colossus/pkg/bucket"
	"github.com/haunt98/colossus/pkg/cache"
	"github.com/haunt98/colossus/pkg/queue"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
		fx.Provide(cachefx.InjectCache(project)),
		fx.Provide(queuefx.InjectQueue(project)),
		fx.Provide(bucketfx.InjectBucket("storage")),
		fx.Provide(grpcfx.InjectGPRCServer(project)),
		fx.Invoke(register),
	).Run()
}

type params struct {
	fx.In

	Sugar  *zap.SugaredLogger
	Cache  *cache.Cache
	Queue  *queue.Queue
	Bucket *bucket.Bucket
	Server *grpc.Server
}

func register(p params) {
	service := ai.NewService(p.Cache, p.Queue, p.Bucket)
	handler := ai.NewHandler(p.Sugar, service)
	handler.Register(p.Server)
}
