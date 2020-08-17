package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
	gatewayv1 "github.com/haunt98/colossus/api/gateway/v1"
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

func (h *Handler) Register(server *grpc.Server) {
	gatewayv1.RegisterGatewayServiceServer(server, h)
}

func (h *Handler) Ping(ctx context.Context, req *gatewayv1.PingRequest) (*gatewayv1.PingResponse, error) {
	return &gatewayv1.PingResponse{}, nil
}

func (h *Handler) Process(ctx context.Context, req *gatewayv1.ProcessRequest) (*gatewayv1.ProcessResponse, error) {
	h.sugar.Infow("Call process", "request", req)

	data := []byte(req.Data)

	id, err := jsonparser.GetString(data, "id")
	if err != nil {
		err = fmt.Errorf("failed to get string %s: %w", "id", err)

		return &gatewayv1.ProcessResponse{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		}, nil
	}

	eventType, err := jsonparser.GetInt(data, "event_type")
	if err != nil {
		err = fmt.Errorf("failed to get string %s: %w", "event_type", err)

		return &gatewayv1.ProcessResponse{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		}, nil
	}

	processInfo, err := h.service.Process(ctx, id, int(eventType))
	if err != nil {
		return &gatewayv1.ProcessResponse{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		}, nil
	}

	return &gatewayv1.ProcessResponse{
		ReturnCode:    processInfo.StatusInfo.Code,
		ReturnMessage: processInfo.StatusInfo.Message,
		TransId:       processInfo.TransID,
	}, nil
}

func (h *Handler) GetStatus(ctx context.Context, req *gatewayv1.GetStatusRequest) (*gatewayv1.GetStatusResponse, error) {
	h.sugar.Infow("Call get status", "request", req)

	processInfo, err := h.service.GetStatus(ctx, req.TransId)
	if err != nil {
		return &gatewayv1.GetStatusResponse{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		}, nil
	}

	data, err := json.Marshal(struct {
		ID string `json:"id"`
	}{
		ID: processInfo.AIOutputID,
	})
	if err != nil {
		err = fmt.Errorf("json failed to marshal: %w", err)

		return &gatewayv1.GetStatusResponse{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		}, nil
	}

	return &gatewayv1.GetStatusResponse{
		ReturnCode:    processInfo.StatusInfo.Code,
		ReturnMessage: processInfo.StatusInfo.Message,
		Data:          string(data),
	}, nil
}
