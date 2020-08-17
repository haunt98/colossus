package ai

import (
	"context"
	"fmt"

	"github.com/haunt98/colossus/pkg/bucket"
	"github.com/haunt98/colossus/pkg/cache"
	"github.com/haunt98/colossus/pkg/queue"
	"github.com/haunt98/colossus/pkg/status"
	"github.com/rs/xid"
)

type Service struct {
	cache         *cache.Cache
	queue         *queue.Queue
	storageBucket *bucket.Bucket
}

func NewService(
	c *cache.Cache,
	q *queue.Queue,
	storageBucket *bucket.Bucket,
) *Service {
	return &Service{
		cache:         c,
		queue:         q,
		storageBucket: storageBucket,
	}
}

func (s *Service) Process(ctx context.Context, id string) (ProcessInfo, error) {
	guid := xid.New()
	transID := guid.String()

	processInfo := ProcessInfo{
		TransID: transID,
		StatusInfo: status.Status{
			Code: status.ProcessingCode,
		},
		InputID: id,
	}

	if err := s.queue.PublishJSON(processInfo); err != nil {
		return ProcessInfo{}, fmt.Errorf("queue failed to publish json: %w", err)
	}

	if err := s.cache.SetJSON(ctx, transID, processInfo); err != nil {
		return ProcessInfo{}, fmt.Errorf("cache failed to set json: %w", err)
	}

	return processInfo, nil
}

func (s *Service) GetStatus(ctx context.Context, transID string) (ProcessInfo, error) {
	var processInfo ProcessInfo
	if err := s.cache.GetJSON(ctx, transID, &processInfo); err != nil {
		return ProcessInfo{}, fmt.Errorf("cache failed to get json: %w", err)
	}

	return processInfo, nil
}
