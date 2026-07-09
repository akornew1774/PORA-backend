// middleware - пакет, обечпечивающий взаимодействие между модулями приложения
package middleware

import (
	"net/http"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware оловит и обрабатывает все паники на сервере
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(
		func(c *gin.Context, recovered interface{}) {
			logger.Log.Error("Паника на сервере: ", recovered)

			err := errors.ErrorInternal
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    err.Code,
				"message": err.Message,
			})

			c.Abort()
		},
	)
}
