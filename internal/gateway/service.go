package gateway

import (
	aiv1 "colossus/api/ai/v1"
	"colossus/pkg/cache"
	"colossus/pkg/status"
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
	"github.com/rs/xid"
)

type Service struct {
	client  *api.Client
	cache   *cache.Cache
	clients map[int]aiv1.AIServiceClient
}

func NewService(
	client *api.Client,
	c *cache.Cache,
	names map[int]string,
) (*Service, error) {
	for eventType, name := range names {
		fmt.Println(eventType, name)
	}

	s1, err := connect.NewService("gateway", client)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := s1.Close(); err != nil {
			log.Println(err)
		}
	}()

	httpClient := s1.HTTPClient()
	rsp, err := httpClient.Get("http://storage.service.consul/ping")
	log.Println(rsp, err)

	return &Service{
		client:  client,
		cache:   c,
		clients: nil,
	}, nil
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
