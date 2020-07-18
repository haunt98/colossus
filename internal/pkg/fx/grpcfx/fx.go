package grpcfx

import (
	"colossus/pkg/network"
	"context"
	"fmt"
	"net"

	"github.com/buger/jsonparser"

	"github.com/hashicorp/consul/api"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ProvideGRPCServerFn func(lc fx.Lifecycle, sugar *zap.SugaredLogger, kv *api.KV, agent *api.Agent) *grpc.Server

func InjectGPRCServer(project string) ProvideGRPCServerFn {
	return func(lc fx.Lifecycle, sugar *zap.SugaredLogger, kv *api.KV, agent *api.Agent) *grpc.Server {
		conf, err := newConfig(kv, project)
		if err != nil {
			sugar.Fatal(err)
		}

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.port))
		if err != nil {
			sugar.Fatalw("Failed to listen", "port", conf.port, "error", err)
		}

		server := grpc.NewServer()

		ip, err := network.GetIP()
		if err != nil {
			sugar.Fatal(err)
		}

		id := fmt.Sprintf("%s-%s-%d", project, ip, conf.port)

		lc.Append(fx.Hook{
			OnStart: onStart(sugar, server, listener, conf, agent, project, ip, id),
			OnStop:  onStop(sugar, server, agent, id),
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

func onStart(sugar *zap.SugaredLogger, server *grpc.Server, listener net.Listener, conf config,
	agent *api.Agent, project, ip, id string) func(context.Context) error {
	return func(ctx context.Context) error {
		sugar.Infow("Start grpc server", "port", conf.port)

		go func() {
			if err := agent.ServiceRegister(&api.AgentServiceRegistration{
				ID:      id,
				Name:    project,
				Port:    int(conf.port),
				Address: ip,
			}); err != nil {
				sugar.Fatalw("Consul agent failed to register", "error", err)
			}

			if err := server.Serve(listener); err != nil {
				sugar.Fatalw("GRPC server failed to serve", "error", err)
			}
		}()

		return nil
	}
}

func onStop(sugar *zap.SugaredLogger, server *grpc.Server,
	agent *api.Agent, id string) func(context.Context) error {
	return func(ctx context.Context) error {
		sugar.Info("Stop GRPC server")

		if err := agent.ServiceDeregister(id); err != nil {
			return fmt.Errorf("consul agent failed to deregister: %w", err)
		}

		server.GracefulStop()

		return nil
	}
}
