package main

import (
	"github.com/gin-gonic/gin"
	"github.com/haunt98/colossus/internal/pkg/fx/amqpfx"
	"github.com/haunt98/colossus/internal/pkg/fx/bucketfx"
	"github.com/haunt98/colossus/internal/pkg/fx/cachefx"
	"github.com/haunt98/colossus/internal/pkg/fx/consulfx"
	"github.com/haunt98/colossus/internal/pkg/fx/ginfx"
	"github.com/haunt98/colossus/internal/pkg/fx/miniofx"
	"github.com/haunt98/colossus/internal/pkg/fx/redisfx"
	"github.com/haunt98/colossus/internal/pkg/fx/zapfx"
	"github.com/haunt98/colossus/internal/storage"
	"github.com/haunt98/colossus/pkg/bucket"
	"github.com/haunt98/colossus/pkg/cache"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	project := "storage"

	fx.New(
		zapfx.Module,
		consulfx.Module,
		redisfx.Module,
		amqpfx.Module,
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
