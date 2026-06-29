// repositories - пакет с методами для работы с БД
package repositories

import (
	"errors"
	"pora/internal/entities"
	"pora/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FamilyRepo - объект, содержащий методы для работы с сущностью семьи в БД
type FamilyRepo struct {
	db *gorm.DB
}

// NewFamilyRepository создает и возвращает новый объект FamilyRepo
func NewFamilyRepository(db *gorm.DB) ports.FamilyRepository {
	return &FamilyRepo{db: db}
}

// FindByID находит семью по ID без связанный с ней сущностей
func (r *FamilyRepo) FindByID(familyID uuid.UUID) (*entities.Family, error) {
	var family entities.Family

	err := r.db.First(&family, familyID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &family, nil
}

// FindWithMembers находит семью по ID вместе с её участниками
func (r *FamilyRepo) FindWithMembers(familyID uuid.UUID) (*entities.Family, error) {
	var family entities.Family

	err := r.db.
		Preload("Members").
		Preload("Members.User").
		First(&family, familyID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &family, nil
}

// FindWithLists находит семью по ID вместе с её списками продуктов
func (r *FamilyRepo) FindWithLists(familyID uuid.UUID) (*entities.Family, error) {
	var family entities.Family

	err := r.db.
		Preload("Lists").
		First(&family, familyID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &family, nil
}

// FindWithAll находит семью по ID вместе со всеми связанными сущностями
func (r *FamilyRepo) FindWithAll(familyID uuid.UUID) (*entities.Family, error) {
	var family entities.Family

	err := r.db.
		Preload("Members").
		Preload("Members.User").
		Preload("Lists").
		First(&family, familyID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &family, nil
}

// CreateFamily создает новый объект семьи в БД
func (r *FamilyRepo) CreateFamily(family *entities.Family) error {
	return r.db.Create(family).Error
}

// UpdateFamily обновляет существующий объект семьи в БД
func (r *FamilyRepo) UpdateFamily(family *entities.Family) error {
	return r.db.Save(family).Error
}

// DeleteFamily удаляет существующий объект семьи в БД
func (r *FamilyRepo) DeleteFamily(family *entities.Family) error {
	return r.db.Delete(family).Error
}
