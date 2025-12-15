package middleware

import (
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StructuredLogger middleware pentru Gin
func StructuredLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Procesează RequestID (X-Request-ID sau generează unul nou)
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)

		// Execută handler-ul următor
		c.Next()

		// După execuție, calculează datele
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// Creăm atributele logului
		attributes := []any{
			slog.String("request_id", requestID),
			slog.Int("status", statusCode),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
			slog.String("user_agent", c.Request.UserAgent()),
		}

		if raw != "" {
			attributes = append(attributes, slog.String("query", raw))
		}

		// Gestionarea erorilor din Gin (c.Error)
		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				attributes = append(attributes, slog.String("error", e))
			}
			logger.Error("HTTP Request Failed", attributes...)
		} else {
			if statusCode >= 500 {
				logger.Error("HTTP Server Error", attributes...)
			} else if statusCode >= 400 {
				logger.Warn("HTTP Client Error", attributes...)
			} else {
				logger.Info("HTTP Request Success", attributes...)
			}
		}
	}
}
