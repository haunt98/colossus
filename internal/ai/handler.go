package ai

import (
	"context"

	aiv1 "github.com/haunt98/colossus/api/ai/v1"
	"github.com/haunt98/colossus/pkg/status"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Handler struct {
	sugar   *zap.SugaredLogger
	service *Service
}

func NewHandler(
	sugar *zap.SugaredLogger,
	service *Service,
) *Handler {
	return &Handler{
		sugar:   sugar,
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
	s.sugar.Infow("Call process", "request", req)
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
	s.sugar.Infow("Call get status", "request", req)
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
