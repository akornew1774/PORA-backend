// ports - пакет, содержащий все порты (интерфейсы)
package ports

import (
	"pora/internal/entities"

	"github.com/google/uuid"
)

// UserRepository содержит порты для взаимодействия с User в БД
type UserRepository interface {
	// FindByID находит пользователя по его ID
	FindByID(userID uuid.UUID) (*entities.User, error)

	// FindByPhone находит пользователя по номеру телефона
	FindByPhone(phone string) (*entities.User, error)

	// CreateUser создает новый объект пользователя в БД
	CreateUser(user *entities.User) error

	// UpdateUser изменяет данные определенного пользователя
	UpdateUser(user *entities.User) error
}

// OTPRepository содержит порты для взаимодействия с Otp в БД
type OtpRepository interface {
	// FindByPhone находит Otp по номеру телефона
	FindByPhone(phone string) (*entities.Otp, error)

	// CreateOTP создает новый объект Otp в БД
	CreateOtp(otp *entities.Otp) error

	// UpdateOTP изменяет данные для сущности OTP
	UpdateOtp(otp *entities.Otp) error
}

// RefreshTokenRepository содержит порты для
// взаимодействия с refresh-токенами в БД
type RefreshTokenRepository interface {
	// FindByToken находит сущность refresh-токена по его значению
	FindByToken(token string) (*entities.RefreshToken, error)

	// CreateRefreshToken создает новый объект токена в БД
	CreateRefreshToken(refreshToken *entities.RefreshToken) error

	// UpdateRefreshToken изменяет данные определенного refresh-токена
	UpdateRefreshToken(refreshToken *entities.RefreshToken) error

	// DeleteByUserID удаляет все refresh-токены c определенным userID
	DeleteByUserID(userID uuid.UUID) error
}
