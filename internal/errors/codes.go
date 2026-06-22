// errors - пакет с сущностями ошибок
package errors

// Задает код ошибки для AppError
const (
	CodeInvalidInput   = "INVALID_INPUT"
	CodeInvalidPhone   = "INVALID_PHONE"
	CodeOTPExpired     = "OTP_EXPIRED"
	CodeOTPIncorrect   = "OTP_INCORRECT"
	CodeAccessExpired  = "ACCESS_TOKEN_EXPIRED"
	CodeRefreshExpired = "REFRESH_TOKEN_EXPIRED"
	CodeNotFound       = "NOT_FOUND"
	CodeForbidden      = "FORBIDDEN"
	CodeInternalError  = "INTERNAL_ERROR"
)
