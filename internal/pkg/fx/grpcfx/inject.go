package grpcfx

import (
	"context"
	"fmt"
	"net"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/consul/api"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ProvideGRPCServerFn func(lc fx.Lifecycle, sugar *zap.SugaredLogger, kv *api.KV) *grpc.Server

func InjectGPRCServer(project string) ProvideGRPCServerFn {
	return func(lc fx.Lifecycle, sugar *zap.SugaredLogger, kv *api.KV) *grpc.Server {
		conf, err := newConfig(kv, project)
		if err != nil {
			sugar.Fatal(err)
		}

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.port))
		if err != nil {
			sugar.Fatalw("Failed to listen", "port", conf.port, "error", err)
		}

		server := grpc.NewServer()

		lc.Append(fx.Hook{
			OnStart: onStart(sugar, server, listener, conf),
			OnStop:  onStop(sugar, server),
		})

		return server
	}
}

type config struct {
	port int
}

func newConfig(kv *api.KV, project string) (config, error) {
	pair, _, err := kv.Get(project, nil)
	if err != nil {
		return config{}, fmt.Errorf("consul kv failed to get key %s: %w", project, err)
	}

	port, err := jsonparser.GetInt(pair.Value, "grpc", "port")
	if err != nil {
		return config{}, fmt.Errorf("failed to get key %s: %w", "grpc.port", err)
	}

	return config{
		port: int(port),
	}, nil
}

func onStart(sugar *zap.SugaredLogger, server *grpc.Server, listener net.Listener, conf config) func(context.Context) error {
	return func(ctx context.Context) error {
		sugar.Infow("Start grpc server", "port", conf.port)

		go func() {
			if err := server.Serve(listener); err != nil {
				sugar.Fatalw("GRPC server failed to serve", "error", err)
			}
		}()

		return nil
	}
}

func onStop(sugar *zap.SugaredLogger, server *grpc.Server) func(context.Context) error {
	return func(ctx context.Context) error {
		sugar.Info("Stop GRPC server")

		server.GracefulStop()

		return nil
	}
}
