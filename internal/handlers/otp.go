// handlers - пакет, содержащий в себе обработчики Api запросов
package handlers

import (
	"net/http"
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"

	"github.com/gin-gonic/gin"
)

// OTPHandler - объект, содержащий методы для обработки
// Api запросов, связанных с OTP-кодами
type OtpHandler struct {
	otpService ports.OtpService
}

// NewOTPHandler создает и возвращает новый объект OTPHandler
func NewOtpHandler(otpService ports.OtpService) *OtpHandler {
	return &OtpHandler{
		otpService: otpService,
	}
}

// SendOTP обрабатывает запрос отправки OTP кода
func (h *OtpHandler) SendOTP(c *gin.Context) {

	var req requests.PhoneAndEmailRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(errors.ErrorInvalidInput)
		return
	}

	err = h.otpService.SendOtp(req.Phone, req.Email)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("OTP-код успешно отправлен на номер: ", req.Phone)
	c.JSON(http.StatusOK, responses.GenericResponse{})
}

// VerifyOTP обрабатывает запрос подтверждения OTP кода
func (h *OtpHandler) VerifyOTP(c *gin.Context) {

	var req requests.VerifyOtpRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(errors.ErrorInvalidInput)
		return
	}

	response, err := h.otpService.VerifyOtp(req.Phone, req.Email, req.OTP)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("OTP-код успешно подтвержден")
	c.JSON(http.StatusOK, response)
}
