package ai

import (
	aiv1 "colossus/api/ai/v1"
	"colossus/pkg/status"
	"context"

	"google.golang.org/grpc"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (s *Handler) Register(server *grpc.Server) {
	aiv1.RegisterAIServiceServer(server, s)
}

func (s *Handler) Ping(ctx context.Context, req *aiv1.PingRequest) (*aiv1.PingResponse, error) {
	return &aiv1.PingResponse{}, nil
}

func (s *Handler) Process(ctx context.Context, req *aiv1.ProcessRequest) (*aiv1.ProcessResponse, error) {
	processInfo, err := s.service.Process(ctx, req.Id)
	if err != nil {
		return &aiv1.ProcessResponse{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		}, nil
	}

	return &aiv1.ProcessResponse{
		ReturnCode:    processInfo.StatusInfo.Code,
		ReturnMessage: processInfo.StatusInfo.Message,
		TransId:       processInfo.TransID,
	}, nil
}

func (s *Handler) GetStatus(ctx context.Context, req *aiv1.GetStatusRequest) (*aiv1.GetStatusResponse, error) {
	processInfo, err := s.service.GetStatus(ctx, req.TransId)
	if err != nil {
		return &aiv1.GetStatusResponse{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		}, nil
	}

	return &aiv1.GetStatusResponse{
		ReturnCode:    processInfo.StatusInfo.Code,
		ReturnMessage: processInfo.StatusInfo.Message,
		Id:            processInfo.OutputID,
	}, nil
}
