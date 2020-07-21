package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SugarGinMiddleware(sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			sugar.Fatal("Failed to load location: %w", err)
		}
		layout := "2006-01-02 15:04:05"

		startTime := time.Now().In(loc)
		path := ctx.Request.URL.Path
		rawQuery := ctx.Request.URL.RawQuery

		// Process request
		ctx.Next()

		stopTime := time.Now().In(loc)

		sugar.Infow("Gin request",
			"start_time", startTime.Format(layout),
			"stop_time", stopTime.Format(layout),
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
