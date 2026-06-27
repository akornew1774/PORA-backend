// repositories - пакет с методами для работы с БД
package repositories

import (
	"errors"
	"pora/internal/entities"
	"pora/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListRepo - объект, содержащий методы
// для работы с сущностью списка покупок в БД
type ListRepo struct {
	db *gorm.DB
}

// NewListRepository создает и возвращает новый объект ListRepo
func NewListRepository(db *gorm.DB) ports.ListRepository {
	return &ListRepo{db: db}
}

// FindByID находит спискок покупок вместе со всеми товарами
func (r *ListRepo) FindByID(listID uuid.UUID) (*entities.List, error) {
	var list entities.List

	err := r.db.
		Preload("Items").
		First(&list, listID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// CreateList создает новый объект списка продуктов в БД
func (r *ListRepo) CreateList(list *entities.List) error {
	return r.db.Create(list).Error
}

// UpdateList обновляет существующий объект списка продуктов в БД
func (r *ListRepo) UpdateList(list *entities.List) error {
	return r.db.Save(list).Error
}

// DeleteList удаляет существующий объект списка продуктов в БД
func (r *ListRepo) DeleteList(list *entities.List) error {
	return r.db.Delete(list).Error
}
