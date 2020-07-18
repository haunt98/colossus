package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SugarGinMiddleware(sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		path := ctx.Request.URL.Path
		rawQuery := ctx.Request.URL.RawQuery

		// Process request
		ctx.Next()

		stopTime := time.Now()

		sugar.Infow("Gin request",
			"start_time", startTime,
			"stop_time", stopTime,
			"latency", stopTime.Sub(startTime),
			"client_ip", ctx.ClientIP(),
			"method", ctx.Request.Method,
			"path", path,
			"raw_query", rawQuery,
			"status_code", ctx.Writer.Status(),
			"error", ctx.Errors.String(),
			"body_size", ctx.Writer.Size(),
		)
	}
}
