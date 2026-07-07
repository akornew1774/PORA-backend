// errors - пакет с сущностями ошибок
package errors

// Задает все нужные для приложения ошибки.
// Все ошибки являются экземплярами AppError
var (
	ErrorInvalidInput = AppError{
		Code:    CodeInvalidInput,
		Message: "Некоректный формат входных данных",
	}

	ErrorInvalidPhone = AppError{
		Code:    CodeInvalidPhone,
		Message: "Некорректный формат номера телефона",
	}

	ErrorOTPExpired = AppError{
		Code:    CodeOTPExpired,
		Message: "Код подтверждения истёк. Попробуйте снова",
	}

	ErrorOTPIncorrect = AppError{
		Code:    CodeOTPIncorrect,
		Message: "Неверный код подтверждения",
	}

	ErrorAccessExpired = AppError{
		Code:    CodeAccessExpired,
		Message: "Сессия истекла. Выполняется обновление токена",
	}

	ErrorRefreshExpired = AppError{
		Code:    CodeRefreshExpired,
		Message: "Сессия истекла. Требуется повторная авторизация",
	}

	ErrorUserNotFound = AppError{
		Code:    CodeNotFound,
		Message: "Пользователь не найден",
	}

	ErrorFamilyNotFound = AppError{
		Code:    CodeNotFound,
		Message: "Семья не найдена",
	}

	ErrorListNotFound = AppError{
		Code:    CodeNotFound,
		Message: "Список продуктов не найден",
	}

	ErrorItemNotFound = AppError{
		Code:    CodeNotFound,
		Message: "Указанный продукт не найден",
	}

	ErrorForbidden = AppError{
		Code:    CodeForbidden,
		Message: "У пользователя нет разрешений на выполнение данного действия",
	}

	ErrorMemberExists = AppError{
		Code:    CodeConflict,
		Message: "Пользователь уже является членом этой семьи",
	}

	ErrorInternal = AppError{
		Code:    CodeInternalError,
		Message: "Непредвиденная ошибка сервера. Попробуйте позже",
	}
)
