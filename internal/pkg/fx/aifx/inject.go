package aifx

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/consul/api"
	"github.com/haunt98/colossus/internal/ai"
	"github.com/tidwall/gjson"
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

	result, err := jsonparser.GetString(pair.Value, "ai", "result")
	if err != nil {
		return ai.CMDConfig{}, fmt.Errorf("failed to get key %s: %w", "ai.result", err)
	}

	return ai.CMDConfig{
		Job:    job,
		Args:   args,
		Result: result,
	}, nil
}

type ProvideNamesFn func(sugar *zap.SugaredLogger, kv *api.KV) map[string]string

func InjectEventTypes(project string) ProvideNamesFn {
	return func(sugar *zap.SugaredLogger, kv *api.KV) map[string]string {
		names, err := newEventTypes(kv, project)
		if err != nil {
			sugar.Fatal(err)
		}

		return names
	}
}

func newEventTypes(kv *api.KV, project string) (map[string]string, error) {
	pair, _, err := kv.Get(project, nil)
	if err != nil {
		return nil, fmt.Errorf("consul kv failed to get key %s: %w", project, err)
	}

	eventTypes := make(map[string]string)
	gEventTypes := gjson.GetBytes(pair.Value, "event_types")
	for i, name := range gEventTypes.Map() {
		eventType := i
		if err != nil {
			return nil, fmt.Errorf("failed to convert string to int: %w", err)
		}

		eventTypes[eventType] = name.String()
	}

	return eventTypes, nil
}

type ProvideURLsFn func(sugar *zap.SugaredLogger, kv *api.KV) map[string]string

func InjectUrls(project string) ProvideURLsFn {
	return func(sugar *zap.SugaredLogger, kv *api.KV) map[string]string {
		urls, err := newURLs(kv, project)
		if err != nil {
			sugar.Fatal(err)
		}

		return urls
	}
}

func newURLs(kv *api.KV, project string) (map[string]string, error) {
	pair, _, err := kv.Get(project, nil)
	if err != nil {
		return nil, fmt.Errorf("consul kv failed to get key %s: %w", project, err)
	}

	urls := make(map[string]string)
	gURLs := gjson.GetBytes(pair.Value, "urls")
	for name, url := range gURLs.Map() {

		urls[name] = url.String()
	}

	return urls, nil
}
