package aifx

import (
	"colossus/internal/ai"
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/buger/jsonparser"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type ProvideCMDConfigFn func(sugar *zap.SugaredLogger, kv *api.KV) ai.CMDConfig

func InjectCMDConfig(project string) ProvideCMDConfigFn {
	return func(sugar *zap.SugaredLogger, kv *api.KV) ai.CMDConfig {
		cmdConf, err := newCMDConfig(kv, project)
		if err != nil {
			sugar.Fatal(err)
		}

		return cmdConf
	}
}

func newCMDConfig(kv *api.KV, project string) (ai.CMDConfig, error) {
	pair, _, err := kv.Get(project, nil)
	if err != nil {
		return ai.CMDConfig{}, fmt.Errorf("consul kv failed to get key %s: %w", project, err)
	}

	job, err := jsonparser.GetString(pair.Value, "ai", "job")
	if err != nil {
		return ai.CMDConfig{}, fmt.Errorf("failed to get key %s: %w", "ai.job", err)
	}

	args := make([]string, 0, 8)
	gArgs := gjson.GetBytes(pair.Value, "ai.args")
	for _, arg := range gArgs.Array() {
		args = append(args, arg.String())
	}
	if len(args) == 0 {
		return ai.CMDConfig{}, fmt.Errorf("empty ai.args")
	}

	outputPath, err := jsonparser.GetString(pair.Value, "ai", "output_path")
	if err != nil {
		return ai.CMDConfig{}, fmt.Errorf("failed to get key %s: %w", "ai.output_path", err)
	}

	return ai.CMDConfig{
		Job:        job,
		Args:       args,
		OutputPath: outputPath,
	}, nil
}
