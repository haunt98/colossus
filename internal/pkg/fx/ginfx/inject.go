package ginfx

import (
	"context"
	"fmt"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/haunt98/colossus/pkg/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ProvideGinFn func(lc fx.Lifecycle, sugar *zap.SugaredLogger, kv *api.KV) *gin.Engine

func InjectGin(project string) ProvideGinFn {
	return func(lc fx.Lifecycle, sugar *zap.SugaredLogger, kv *api.KV) *gin.Engine {
		conf, err := newConfig(kv, project)
		if err != nil {
			sugar.Fatal(err)
		}

		gin.SetMode(conf.mode)
		engine := gin.New()
		engine.Use(gin.Recovery())
		engine.Use(middleware.SugarGinMiddleware(sugar))

		engine.GET("/ping", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "pong")
		})

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", conf.port),
			Handler: engine,
		}

		lc.Append(fx.Hook{
			OnStart: onStart(sugar, server, conf),
			OnStop:  onStop(sugar, server),
		})

		return engine
	}
}

type config struct {
	mode string
	port int
}

func newConfig(kv *api.KV, project string) (config, error) {
	pair, _, err := kv.Get(project, nil)
	if err != nil {
		return config{}, fmt.Errorf("consul kv failed to get key %s: %w", project, err)
	}

	mode, err := jsonparser.GetString(pair.Value, "gin", "mode")
	if err != nil {
		return config{}, fmt.Errorf("failed to get key %s: %w", "gin.mode", err)
	}

	port, err := jsonparser.GetInt(pair.Value, "gin", "port")
	if err != nil {
		return config{}, fmt.Errorf("failed to get key %s: %w", "gin.port", err)
	}

	return config{
		mode: mode,
		port: int(port),
	}, nil
}

func onStart(sugar *zap.SugaredLogger, server *http.Server, conf config) func(context.Context) error {
	return func(ctx context.Context) error {
		sugar.Infow("Start gin", "mode", conf.mode, "port", conf.port)

		go func() {
			if err := server.ListenAndServe(); err != nil {
				sugar.Fatalw("HTTP server failed to listen and serve", "error", err)
			}
		}()

		return nil
	}
}

func onStop(sugar *zap.SugaredLogger, server *http.Server) func(context.Context) error {
	return func(ctx context.Context) error {
		sugar.Info("Stop gin")

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("http server failed to shutdown: %w", err)
		}

		return nil
	}
}
