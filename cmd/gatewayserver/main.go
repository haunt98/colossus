package main

import (
	"colossus/internal/gateway"
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
		fx.Invoke(register),
	).Run()
}

type params struct {
	fx.In

	Cache  *cache.Cache
	Agent  *api.Agent
	Server *grpc.Server
}

func register(p params) {
	service := gateway.NewService(p.Cache, p.Agent, nil)
	handler := gateway.NewHandler(service)
	handler.Register(p.Server)
}
