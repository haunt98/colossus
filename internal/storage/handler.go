package storage

import (
	"colossus/pkg/status"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
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

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusOK, Response{
			ReturnCode:    status.FailedCode,
			ReturnMessage: err.Error(),
		})
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			h.sugar.Error(err)
		}
	}()

	fileInfo, err := h.service.Upload(ctx, file, fileHeader.Size)
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
