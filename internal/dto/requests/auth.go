// requests - пакет, содержащий структуры запросов по Api
package requests

// VerifyOTPRequest - структура для запроса на подтверждение OTP кода
type VerifyOtpRequest struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
	OTP   string `json:"otp" binding:"required"`
}

// RefreshRequest - структура для запроса на обновление refresh токена
type RefreshRequest struct {
	RefreshToken string `form:"refresh-token" binding:"required"`
}
