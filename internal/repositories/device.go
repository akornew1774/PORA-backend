// repositories - пакет с методами для работы с БД
package repositories

import (
	"errors"
	"pora/internal/entities"
	"pora/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceRepo - объект, содержащий методы для
// работы с сущностью устройства пользователя в БД
type DeviceRepo struct {
	db *gorm.DB
}

// NewDeviceRepo создает и возвращает новый объект DeviceRepo
func NewDeviceRepository(db *gorm.DB) ports.DeviceRepository {
	return &DeviceRepo{db: db}
}

// FindByUserID находит устройство пользователя по ID пользователя
func (r *DeviceRepo) FindByUserID(
	userID uuid.UUID) (*entities.Device, error) {

	var device entities.Device

	err := r.db.
		Where("user_id = ?", userID).
		First(&device).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// CreateDevice создает новый объект устройства в БД
func (r *DeviceRepo) CreateDevice(device *entities.Device) error {
	return r.db.Create(device).Error
}

// UpdateDevice обновляет существующий объект устройства в БД
func (r *DeviceRepo) UpdateDevice(device *entities.Device) error {
	return r.db.Save(device).Error
}

// DeleteDevice удаляет существующий объект устройства в БД
func (r *DeviceRepo) DeleteDevice(device *entities.Device) error {
	return r.db.Delete(device).Error
}
