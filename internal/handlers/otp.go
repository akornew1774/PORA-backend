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
type OTPHandler struct {
	otpService ports.OtpService
}

// NewOTPHandler создает и возвращает новый объект OTPHandler
func NewOTPHandler(otpService ports.OtpService) *OTPHandler {
	return &OTPHandler{
		otpService: otpService,
	}
}

// SendOTP обрабатывает запрос отправки OTP кода
func (h *OTPHandler) SendOTP(c *gin.Context) {
	var req requests.PhoneNumberRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(errors.ErrorInvalidInput)
		return
	}

	err = h.otpService.SendOtp(req.Phone)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("OTP-код успешно отправлен на номер: ", req.Phone)
	c.JSON(http.StatusOK, responses.GenericResponse{})
}

// VerifyOTP обрабатывает запрос подтверждения OTP кода
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req requests.VerifyOtpRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(errors.ErrorInvalidInput)
		return
	}

	response, err := h.otpService.VerifyOtp(req.Phone, req.OTP)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("OTP-код успешно подтвержден")
	c.JSON(http.StatusOK, response)
}
