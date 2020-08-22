package main

import (
	"github.com/haunt98/colossus/internal/gateway"
	"github.com/haunt98/colossus/internal/pkg/fx/aifx"
	"github.com/haunt98/colossus/internal/pkg/fx/amqpfx"
	"github.com/haunt98/colossus/internal/pkg/fx/bucketfx"
	"github.com/haunt98/colossus/internal/pkg/fx/cachefx"
	"github.com/haunt98/colossus/internal/pkg/fx/consulfx"
	"github.com/haunt98/colossus/internal/pkg/fx/grpcfx"
	"github.com/haunt98/colossus/internal/pkg/fx/miniofx"
	"github.com/haunt98/colossus/internal/pkg/fx/queuefx"
	"github.com/haunt98/colossus/internal/pkg/fx/redisfx"
	"github.com/haunt98/colossus/internal/pkg/fx/zapfx"
	"github.com/haunt98/colossus/pkg/cache"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	project := "gateway"

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
		fx.Provide(fx.Annotated{
			Name:   "event_types",
			Target: aifx.InjectEventTypes(project),
		}),
		fx.Provide(fx.Annotated{
			Name:   "urls",
			Target: aifx.InjectUrls(project),
		}),
		fx.Invoke(register),
	).Run()
}

type params struct {
	fx.In

	Sugar      *zap.SugaredLogger
	Cache      *cache.Cache
	Server     *grpc.Server
	EventTypes map[string]string `name:"event_types"`
	URLs       map[string]string `name:"urls"`
}

func register(p params) {
	service, err := gateway.NewService(p.Sugar, p.Cache, p.EventTypes, p.URLs)
	if err != nil {
		p.Sugar.Fatal(err)
	}

	handler := gateway.NewHandler(p.Sugar, service)
	handler.Register(p.Server)
}
