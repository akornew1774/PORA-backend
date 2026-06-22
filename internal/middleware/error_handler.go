// middleware - пакет, обечпечивающий взаимодействие между модулями приложения
package middleware

import (
	"net/http"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

// ErrorHandler обрабатывает ошибки, полученные в ходе обработки запроса.
// Выполняется после выполнения хэндлера, обрабатывает последнюю из полученных ошибок
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		if appErr, ok := err.(errors.AppError); ok {

			logger.Log.Warn("Возникла ошибка при обработке запроса: ", appErr.Code)

			switch appErr.Code {

			case errors.CodeInvalidInput,
				errors.CodeInvalidPhone,
				errors.CodeOTPExpired,
				errors.CodeOTPIncorrect:

				c.JSON(http.StatusBadRequest, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				return

			case errors.CodeAccessExpired,
				errors.CodeRefreshExpired:

				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				return

			case errors.CodeNotFound:

				c.JSON(http.StatusNotFound, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				return

			case errors.CodeForbidden:

				c.JSON(http.StatusForbidden, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				return

			default:

				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    errors.ErrorInternal.Code,
					"message": errors.ErrorInternal.Message,
				})
				return
			}
		}

		logger.Log.Error("Непредвиденная ошибка на сервере: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrorInternal.Code,
			"message": errors.ErrorInternal.Message,
		})
	}
}
