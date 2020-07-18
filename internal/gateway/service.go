package gateway

import (
	aiv1 "colossus/api/ai/v1"
	"colossus/pkg/cache"
	"colossus/pkg/status"
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"

	"github.com/rs/xid"
)

type Service struct {
	cache   *cache.Cache
	agent   *api.Agent
	clients map[int]aiv1.AIServiceClient
}

func NewService(
	c *cache.Cache,
	agent *api.Agent,
	names map[int]string,
) *Service {
	for eventType, name := range names {
		fmt.Println(eventType, name)
	}

	return &Service{
		cache:   c,
		clients: nil,
	}
}

func (s *Service) Process(ctx context.Context, id string, eventType int) (ProcessInfo, error) {
	client, ok := s.clients[eventType]
	if !ok {
		return ProcessInfo{}, fmt.Errorf("event_type %d is unknown", eventType)
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
		return ProcessInfo{}, fmt.Errorf("event_type %d is unknown", processInfo.EventType)
	}

	rsp, err := client.GetStatus(ctx, &aiv1.GetStatusRequest{
		TransId: processInfo.TransID,
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
