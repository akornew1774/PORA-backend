// repositories - пакет с методами для работы с БД
package repositories

import (
	"errors"
	"pora/internal/entities"
	"pora/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ItemRepo - объект, содержащий методы для работы с сущностью товара в БД
type ItemRepo struct {
	db *gorm.DB
}

// NewItemRepository создает и возвращает новый объект ItemRepo
func NewItemRepository(db *gorm.DB) ports.ItemRepository {
	return &ItemRepo{db: db}
}

// FindByID находит товар из списка продуктов по ID
func (r *ItemRepo) FindByID(itemID uuid.UUID) (*entities.Item, error) {
	var item entities.Item

	err := r.db.First(&item, itemID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// CreateItem создает новый объект товара в БД
func (r *ItemRepo) CreateItem(item *entities.Item) error {
	return r.db.Create(item).Error
}

// UpdateItem обновляет существующий объект товара в БД
func (r *ItemRepo) UpdateItem(item *entities.Item) error {
	return r.db.Save(item).Error
}

// DeleteItem удаляет существующий объект товара в БД
func (r *ItemRepo) DeleteItem(item *entities.Item) error {
	return r.db.Delete(item).Error
}
