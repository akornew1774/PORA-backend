// ports - пакет, содержащий все порты (интерфейсы)
package ports

import (
	"io"
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
	"pora/internal/entities"

	"github.com/google/uuid"
)

// UserService содержит порты для методов, необходимых для
// работы с пользователем и его профилем
type UserService interface {
	// UpdateUserInfo обновляет (или дополняет) информацию
	// о каком-то конкретном пользователе
	UpdateUserInfo(userID uuid.UUID,
		req *requests.UpdateUserRequest) error

	// SaveImage вызывает FileService для сохранения изображения в
	// хранилище, сохраняет ссылку на него в профиле пользователя в БД
	SaveImage(userID uuid.UUID,
		file io.Reader, fileName string) (string, error)
}

// OtpService содержит порты для методов отправки и подтверждения OTP
type OtpService interface {
	// SendOtp отправляет Otp-код по номеру телефона и сохраняет его
	SendOtp(rawPhone string) error

	// VerifyOtp сравнивает полученный Otp-код c сохраненным в БД
	VerifyOtp(rawPhone string, otp string) (
		responses.VerifyOtpResponse, error)
}

// AuthService содержит порты для методов, необходимых для авторизации
type AuthService interface {
	// CheckUser находит пользователя по телефону и проверяет его
	// статус (notFound / notRegistered / registered)
	CheckUser(rawPhone string) (responses.IsUserResponse, error)

	// GetNewTokens получает новые access и refresh токены для пользователя
	GetNewTokens(refreshTokenStr string) (responses.RefreshResponse, error)

	// RemoveRefreshToken удаляет refresh-токен по ID пользователя
	RemoveRefreshToken(userID uuid.UUID) error
}

// TokenService содержит порты для методов создания/обновления токенов
type TokenService interface {
	// CreateTokens создает новый access-токен,
	// обновляет (или создает новый) refresh-токен,
	// после чего возвращает access и refresh токены.
	CreateTokens(user *entities.User) (string, string, error)

	// DecodeAccessToken декодирует полученный access-токен,
	// проверяет его валидность
	DecodeAccessToken(rawToken string) (uuid.UUID, error)
}

// FileService содержит порты для методов для
// работы с изображениями в хранилище
type FileService interface {
	// SaveProfileImage сохраняет изображение профиля пользователя
	// в хранилище и возвращает ссылку на него
	SaveProfileImage(file io.Reader,
		fileName string, userID uuid.UUID) (string, error)
}

// FamilyService содержит порты для методов для
// создания семей и просмотра информации о них
type FamilyService interface {
}

// ListService содержит порты для методов для
// просмотра и изменения списков покупок
type ListService interface {
	// GetAllSections делит все товары на секции и
	// переводит их в нужный для response формат
	GetAllSections(list *entities.List) (
		[]responses.SectionInfo, error)

	// GetHighestPrioritySections отбирает товары высшего приоритета, делит
	// их на секции и приводит к нужному для response формату
	GetHighestPrioritySections(list *entities.List) (
		[]responses.SectionInfo, error)
}

// ItemService содержит порты для методов для
// изменения и удаления продуктов, а также прочих
// взаимодействий с данной сущностью
type ItemService interface {
}
