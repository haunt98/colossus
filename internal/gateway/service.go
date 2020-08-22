package gateway

import (
	"context"
	"fmt"

	aiv1 "github.com/haunt98/colossus/api/ai/v1"
	"github.com/haunt98/colossus/pkg/cache"
	"github.com/haunt98/colossus/pkg/status"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Service struct {
	cache   *cache.Cache
	clients map[string]aiv1.AIServiceClient
}

func NewService(
	sugar *zap.SugaredLogger,
	c *cache.Cache,
	eventTypes map[string]string,
	urls map[string]string,
) (*Service, error) {
	sugar.Infow("Init service", "event_types", eventTypes, "urls", urls)

	clients := make(map[string]aiv1.AIServiceClient, len(urls))
	for eventType, name := range eventTypes {
		url, ok := urls[name]
		if !ok {
			sugar.Error("url unknown", "name", name)
			continue
		}

		conn, err := grpc.Dial(url, grpc.WithInsecure())
		if err != nil {
			sugar.Errorw("GRPC failed to dial", "error", err)
			continue
		}

		client := aiv1.NewAIServiceClient(conn)
		clients[eventType] = client
	}

	return &Service{
		cache:   c,
		clients: clients,
	}, nil
}

func (s *Service) Process(ctx context.Context, id string, eventType string) (ProcessInfo, error) {
	client, ok := s.clients[eventType]
	if !ok {
		return ProcessInfo{}, fmt.Errorf("eventType %s is unknown", eventType)
	}

	rsp, err := client.Process(ctx, &aiv1.ProcessRequest{
		Id: id,
	})
	if err != nil {
		return ProcessInfo{}, fmt.Errorf("client failed to process: %w", err)
	}

	guid := xid.New()
	transID := guid.String()

	processInfo := ProcessInfo{
		TransID: transID,
		StatusInfo: status.Status{
			Code:    rsp.ReturnCode,
			Message: rsp.ReturnMessage,
		},
		EventType: eventType,
		AITransID: rsp.TransId,
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

	client, ok := s.clients[processInfo.EventType]
	if !ok {
		return ProcessInfo{}, fmt.Errorf("eventType %s is unknown", processInfo.EventType)
	}

	rsp, err := client.GetStatus(ctx, &aiv1.GetStatusRequest{
		TransId: processInfo.AITransID,
	})
	if err != nil {
		return ProcessInfo{}, fmt.Errorf("client failed to get status: %w", err)
	}

	processInfo.StatusInfo = status.Status{
		Code:    rsp.ReturnCode,
		Message: rsp.ReturnMessage,
	}
	processInfo.AIOutputID = rsp.Id

	if err := s.cache.SetJSON(ctx, transID, processInfo); err != nil {
		return ProcessInfo{}, fmt.Errorf("cache failed to set json: %w", err)
	}

	return processInfo, nil
}
