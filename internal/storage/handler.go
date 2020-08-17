package storage

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haunt98/colossus/pkg/status"
	"go.uber.org/zap"
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

func (h *Handler) Register(engine *gin.Engine) {
	group := engine.Group("")
	{
		group.POST("/upload", h.upload)
		group.GET("/download", h.download)
	}
}

func (h *Handler) upload(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, Response{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		})
		return
	}

	fileInfo, err := h.service.Upload(ctx, h.sugar, fileHeader)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		ReturnCode: status.SuccessfulCode,
		Data:       fileInfo,
	})
}

func (h *Handler) download(ctx *gin.Context) {
	id := ctx.Query("id")

	url, err := h.service.Download(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		ReturnCode: status.SuccessfulCode,
		Data:       url,
	})
}
