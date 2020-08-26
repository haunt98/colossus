package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/haunt98/colossus/internal/storage"
	"github.com/haunt98/colossus/pkg/bucket"
	"github.com/haunt98/colossus/pkg/cache"
	"github.com/haunt98/colossus/pkg/queue"
	"github.com/haunt98/colossus/pkg/status"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type Processor struct {
	sugar          *zap.SugaredLogger
	cache          *cache.Cache
	queue          *queue.Queue
	storageCache   *cache.Cache
	strorageBucket *bucket.Bucket
	cmdConf        CMDConfig
}

func NewProcessor(
	sugar *zap.SugaredLogger,
	c *cache.Cache,
	q *queue.Queue,
	storageCache *cache.Cache,
	storageBucket *bucket.Bucket,
	cmdConf CMDConfig,
) *Processor {
	sugar.Infow("Init processor", "cmd_config", cmdConf)

	return &Processor{
		sugar:          sugar,
		cache:          c,
		queue:          q,
		storageCache:   storageCache,
		strorageBucket: storageBucket,
		cmdConf:        cmdConf,
	}
}

func (p *Processor) Consume() error {
	deliveries, err := p.queue.Consume()
	if err != nil {
		return err
	}

	for delivery := range deliveries {
		if err := p.consumeBody(delivery.Body); err != nil {
			p.sugar.Errorw("Failed to consume body", "error", err)
		}
	}

	return nil
}

func (p *Processor) consumeBody(body []byte) error {
	var processInfo ProcessInfo
	if err := json.Unmarshal(body, &processInfo); err != nil {
		return fmt.Errorf("json failed to unmarshal: %w", err)
	}

	p.sugar.Infow("Consume", "process_info", processInfo)

	if err := p.process(&processInfo); err != nil {
		p.sugar.Errorw("Failed to process", "error", err)
		processInfo.StatusInfo = status.Status{
			Code:    status.FailedCode,
			Message: err.Error(),
		}
	}

	p.sugar.Infow("After", "process_info", processInfo)

	if err := p.cache.SetJSON(context.Background(), processInfo.TransID, processInfo); err != nil {
		return fmt.Errorf("cache failed to set json: %w", err)
	}

	return nil
}

const tmpDir = "tmp/"

func (p *Processor) process(processInfo *ProcessInfo) error {
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return fmt.Errorf("os failed to mkdir all: %w", err)
	}

	inputPath, inputPathWithExt, err := p.downloadInput(processInfo, tmpDir)
	if err != nil {
		return fmt.Errorf("failed to download input: %w", err)
	}

	guid := xid.New()
	outputID := guid.String()
	outputPath := tmpDir + outputID

	cmdConf := p.cmdConf.transform(processInfo.InputID, inputPath, inputPathWithExt, outputID, outputPath)
	p.sugar.Infow("Actual", "cmdConfig", cmdConf)
	cmdOutput, err := exec.Command(cmdConf.Job, cmdConf.Args...).Output()
	if err != nil {
		log.Printf("failed to run: %s\n", err)
	}

	cmdOutputPathWithExt := outputPath + "_cmd_output.txt"
	if err := ioutil.WriteFile(cmdOutputPathWithExt, cmdOutput, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", cmdOutputPathWithExt, err)
	}

	outputPath = cmdConf.Result
	p.sugar.Infow("Actual", "outputPath", outputPath)

	if err := p.uploadOutput(processInfo, outputID, outputPath); err != nil {
		return fmt.Errorf("failed to upload output: %w", err)
	}

	processInfo.StatusInfo = status.Status{
		Code: status.SuccessfulCode,
	}
	processInfo.OutputID = outputID

	return nil
}

func (p *Processor) downloadInput(processInfo *ProcessInfo, tmpDir string) (string, string, error) {
	var inputFileInfo storage.FileInfo
	if err := p.storageCache.GetJSON(context.Background(), processInfo.InputID, &inputFileInfo); err != nil {
		return "", "", fmt.Errorf("cache failed to get json: %w", err)
	}

	inputPathWithExt := tmpDir + inputFileInfo.ID + inputFileInfo.Extension
	if err := p.strorageBucket.FGetObject(processInfo.InputID, inputPathWithExt); err != nil {
		return "", "", fmt.Errorf("bucket failed to fget object: %w", err)
	}
	inputPath := tmpDir + inputFileInfo.ID

	return inputPath, inputPathWithExt, nil
}

func (p *Processor) uploadOutput(processInfo *ProcessInfo, outputID, outputPath string) error {
	outputContentType, err := mimetype.DetectFile(outputPath)
	if err != nil {
		return fmt.Errorf("mimetype failed to detect file: %w", err)
	}

	if err := p.strorageBucket.FPutObject(outputID, outputPath, outputContentType.String()); err != nil {
		return fmt.Errorf("bucket failed to fput object: %w", err)
	}

	var outputFileInfo = storage.FileInfo{
		ID:          outputID,
		ContentType: outputContentType.String(),
		Extension:   outputContentType.Extension(),
	}

	if err := p.storageCache.SetJSON(context.Background(), outputFileInfo.ID, outputFileInfo); err != nil {
		return fmt.Errorf("cache failed to set json: %w", err)
	}

	return nil
}

type CMDConfig struct {
	Job    string
	Args   []string
	Result string
}

func (conf CMDConfig) transform(inputID, inputPath, inputPathWithExt,
	outputID, result string) CMDConfig {
	r := strings.NewReplacer("{input_id}", inputID, "{input_path}", inputPath,
		"{input_path_with_ext}", inputPathWithExt,
		"{output_id}", outputID, "{output_path}", result)

	args := make([]string, len(conf.Args))
	for i := range conf.Args {
		args[i] = r.Replace(conf.Args[i])
	}

	result = r.Replace(conf.Result)

	return CMDConfig{
		Job:    conf.Job,
		Args:   args,
		Result: result,
	}
}
