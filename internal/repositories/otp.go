// repositories - пакет с методами для работы с БД
package repositories

import (
	"errors"
	"pora/internal/entities"
	"pora/internal/ports"

	"gorm.io/gorm"
)

// OtpRepo - объект, содержащий методы для работы с сущностью Otp кода
type OtpRepo struct {
	db *gorm.DB
}

// NewOTPRepository создает и возвращает новый объект otpRepo
func NewOtpRepository(db *gorm.DB) ports.OtpRepository {
	return &OtpRepo{db: db}
}

// FindByPhone находит Otp по номеру телефона
func (r *OtpRepo) FindByPhone(phone string) (*entities.Otp, error) {
	var otp entities.Otp

	err := r.db.Where("phone = ?", phone).
		Order("expires_at DESC").
		First(&otp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// CreateOTP создает новый объект Otp в БД
func (r *OtpRepo) CreateOtp(otp *entities.Otp) error {
	return r.db.Create(otp).Error
}

// UpdateOTP изменяет данные для сущности OTP
func (r *OtpRepo) UpdateOtp(otp *entities.Otp) error {
	return r.db.Save(otp).Error
}
