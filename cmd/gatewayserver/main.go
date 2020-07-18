package main

import (
	"colossus/internal/gateway"
	"colossus/internal/pkg/fx/aifx"
	"colossus/internal/pkg/fx/bucketfx"
	"colossus/internal/pkg/fx/cachefx"
	"colossus/internal/pkg/fx/consulfx"
	"colossus/internal/pkg/fx/grpcfx"
	"colossus/internal/pkg/fx/miniofx"
	"colossus/internal/pkg/fx/queuefx"
	"colossus/internal/pkg/fx/rabbitmqfx"
	"colossus/internal/pkg/fx/redisfx"
	"colossus/internal/pkg/fx/zapfx"
	"colossus/pkg/cache"

	"go.uber.org/zap"

	"github.com/hashicorp/consul/api"

	"google.golang.org/grpc"

	"go.uber.org/fx"
)

func main() {
	project := "gateway"

	fx.New(
		zapfx.Module,
		consulfx.Module,
		redisfx.Module,
		rabbitmqfx.Module,
		miniofx.Module,
		fx.Provide(cachefx.InjectCache(project)),
		fx.Provide(queuefx.InjectQueue(project)),
		fx.Provide(bucketfx.InjectBucket("storage")),
		fx.Provide(grpcfx.InjectGPRCServer(project)),
		fx.Provide(fx.Annotated{
			Name:   "names",
			Target: aifx.InjectNames(project),
		}),
		fx.Invoke(register),
	).Run()
}

type params struct {
	fx.In

	Sugar  *zap.SugaredLogger
	Client *api.Client
	Cache  *cache.Cache
	Server *grpc.Server
	Names  map[int]string `name:"names"`
}

func register(p params) {
	service, err := gateway.NewService(p.Client, p.Cache, p.Names)
	if err != nil {
		p.Sugar.Fatal(err)
	}

	handler := gateway.NewHandler(service)
	handler.Register(p.Server)
}
