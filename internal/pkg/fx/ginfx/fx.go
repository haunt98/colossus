package ginfx

import (
	"colossus/pkg/middleware"
	"colossus/pkg/network"
	"context"
	"fmt"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ProvideGinFn func(lc fx.Lifecycle, sugar *zap.SugaredLogger, kv *api.KV, agent *api.Agent) *gin.Engine

func InjectGin(project string) ProvideGinFn {
	return func(lc fx.Lifecycle, sugar *zap.SugaredLogger, kv *api.KV, agent *api.Agent) *gin.Engine {
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

		ip, err := network.GetIP()
		if err != nil {
			sugar.Fatal(err)
		}

		id := fmt.Sprintf("%s-%s-%d", project, ip, conf.port)

		lc.Append(fx.Hook{
			OnStart: onStart(sugar, server, conf, agent, project, ip, id),
			OnStop:  onStop(sugar, server, agent, id),
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

func onStart(sugar *zap.SugaredLogger, server *http.Server, conf config,
	agent *api.Agent, project, ip, id string) func(context.Context) error {
	return func(ctx context.Context) error {
		sugar.Infow("Start gin", "mode", conf.mode, "port", conf.port)

		go func() {
			if err := agent.ServiceRegister(&api.AgentServiceRegistration{
				ID:      id,
				Name:    project,
				Port:    int(conf.port),
				Address: ip,
			}); err != nil {
				sugar.Fatalw("Consul agent failed to register", "error", err)
			}

			if err := server.ListenAndServe(); err != nil {
				sugar.Fatalw("HTTP server failed to listen and serve", "error", err)
			}
		}()

		return nil
	}
}

func onStop(sugar *zap.SugaredLogger, server *http.Server,
	agent *api.Agent, id string) func(context.Context) error {
	return func(ctx context.Context) error {
		sugar.Info("Stop gin")

		if err := agent.ServiceDeregister(id); err != nil {
			return fmt.Errorf("consul agent failed to deregister: %w", err)
		}

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("http server failed to shutdown: %w", err)
		}

		return nil
	}
}
