package zapfx

import (
	"context"
	"log"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(
	provideLogger,
	provideSugar,
)

func provideLogger(lc fx.Lifecycle) *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err.Error())
	}

	lc.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(ctx context.Context) error {
			return logger.Sync()
		},
	})

	return logger
}

func provideSugar(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}
