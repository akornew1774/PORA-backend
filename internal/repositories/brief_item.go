// repositories - пакет с методами для работы с БД
package repositories

import (
	"pora/internal/entities"
	"pora/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BriefItemRepo - объект, содержащий методы для
// работы с краткой информацией о товаре в БД
type BriefItemRepo struct {
	db *gorm.DB
}

// NewBriefItemRepository возвращает новый объект BriefItemRepo
func NewBriefItemRepository(db *gorm.DB) ports.BriefItemRepository {
	return &BriefItemRepo{db: db}
}

// FindByUserID удаляет все объекты BriefItem,
// привязанные к определенному пользователю
func (r *BriefItemRepo) FindByUserID(
	userID uuid.UUID) ([]entities.BriefItem, error) {

	var briefItems []entities.BriefItem

	err := r.db.
		Where("user_id = ?", userID).
		Find(&briefItems).
		Error

	return briefItems, err
}

// CreateBriefItem создает новый объект BriefItem в БД
func (r *BriefItemRepo) CreateBriefItem(
	briefItem *entities.BriefItem) error {

	return r.db.Create(briefItem).Error
}

// DeleteByUserID удаляет все объекты BriefItem,
// привязанные к определенному пользователю
func (r *BriefItemRepo) DeleteByUserID(
	userID uuid.UUID) error {

	err := r.db.
		Where("user_id = ?", userID).
		Delete(&entities.BriefItem{}).
		Error

	return err
}
