package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		fields := logrus.Fields{
			"client_ip":    param.ClientIP,
			"method":       param.Method,
			"path":         param.Path,
			"status_code":  param.StatusCode,
			"latency":      param.Latency,
			"latency_ms":   float64(param.Latency.Nanoseconds()) / 1000000.0,
			"user_agent":   param.Request.UserAgent(),
			"request_size": param.Request.ContentLength,
			"timestamp":    param.TimeStamp.Format(time.RFC3339),
		}

		if isPasswordEndpoint(param.Method, param.Path) {
			fields["sensitive_data"] = true
			fields["note"] = "request body hidden for security"
		}

		logger.WithFields(fields).Info("HTTP Request")
		return ""
	})
}

func isPasswordEndpoint(method, path string) bool {
	return method == "POST" && strings.Contains(path, "/users")
}

func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := logrus.Fields{
			"request_id": c.Request.Header.Get("X-Request-ID"),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
		}

		requestLogger := logger.WithFields(fields)
		c.Set("logger", requestLogger)
		c.Next()
	}
}
