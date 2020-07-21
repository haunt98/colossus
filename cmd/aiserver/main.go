package main

import (
	"colossus/internal/ai"
	"colossus/internal/pkg/fx/amqpfx"
	"colossus/internal/pkg/fx/bucketfx"
	"colossus/internal/pkg/fx/cachefx"
	"colossus/internal/pkg/fx/consulfx"
	"colossus/internal/pkg/fx/grpcfx"
	"colossus/internal/pkg/fx/miniofx"
	"colossus/internal/pkg/fx/queuefx"
	"colossus/internal/pkg/fx/redisfx"
	"colossus/internal/pkg/fx/zapfx"
	"colossus/pkg/bucket"
	"colossus/pkg/cache"
	"colossus/pkg/queue"
	"log"
	"os"

	"google.golang.org/grpc"

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
	handler := ai.NewHandler(service)
	handler.Register(p.Server)
}
