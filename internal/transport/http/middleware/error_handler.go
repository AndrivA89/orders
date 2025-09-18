package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func ErrorHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.WithFields(logrus.Fields{
					"panic":  err,
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
				}).Error("Panic recovered")

				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Internal Server Error",
					Message: "An unexpected error occurred",
					Code:    "INTERNAL_ERROR",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}

func HandleError(c *gin.Context, statusCode int, err error, code string) {
	logger, exists := c.Get("logger")
	if exists {
		if requestLogger, ok := logger.(*logrus.Entry); ok {
			requestLogger.WithFields(logrus.Fields{
				"error":       err.Error(),
				"status_code": statusCode,
				"error_code":  code,
			}).Error("Request error")
		}
	}

	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: err.Error(),
		Code:    code,
	}

	c.JSON(statusCode, response)
}

func HandleValidationError(c *gin.Context, err error) {
	HandleError(c, http.StatusBadRequest, err, "VALIDATION_ERROR")
}

func HandleNotFoundError(c *gin.Context, err error) {
	HandleError(c, http.StatusNotFound, err, "NOT_FOUND")
}

func HandleInternalError(c *gin.Context, err error) {
	HandleError(c, http.StatusInternalServerError, err, "INTERNAL_ERROR")
}
