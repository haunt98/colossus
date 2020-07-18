package main

import (
	"colossus/internal/pkg/fx/bucketfx"
	"colossus/internal/pkg/fx/cachefx"
	"colossus/internal/pkg/fx/consulfx"
	"colossus/internal/pkg/fx/ginfx"
	"colossus/internal/pkg/fx/miniofx"
	"colossus/internal/pkg/fx/rabbitmqfx"
	"colossus/internal/pkg/fx/redisfx"
	"colossus/internal/pkg/fx/zapfx"
	"colossus/internal/storage"
	"colossus/pkg/bucket"
	"colossus/pkg/cache"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go.uber.org/fx"
)

func main() {
	project := "storage"

	fx.New(
		zapfx.Module,
		consulfx.Module,
		redisfx.Module,
		rabbitmqfx.Module,
		miniofx.Module,
		fx.Provide(cachefx.InjectCache(project)),
		fx.Provide(bucketfx.InjectBucket(project)),
		fx.Provide(ginfx.InjectGin(project)),
		fx.Invoke(register),
	).Run()
}

type params struct {
	fx.In

	Sugar  *zap.SugaredLogger
	Bucket *bucket.Bucket
	Cache  *cache.Cache
	Engine *gin.Engine
}

func register(p params) {
	service := storage.NewService(p.Bucket, p.Cache)
	handler := storage.NewHandler(p.Sugar, service)
	handler.Register(p.Engine)
}
