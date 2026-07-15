// repositories - пакет с методами для работы с БД
package repositories

import (
	"errors"
	"pora/internal/entities"
	"pora/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MemberRepo - объект, содержащий методы
// для работы с участниками семьи
type MemberRepo struct {
	db *gorm.DB
}

// NewMemberRepository создает и возвращает новый объект MemberRepo
func NewMemberRepository(db *gorm.DB) ports.MemberRepository {
	return &MemberRepo{db: db}
}

// FindByUserAndFamily находит участника по его ID и ID семьи
func (r *MemberRepo) FindByUserAndFamily(userID uuid.UUID,
	familyID uuid.UUID) (*entities.FamilyMember, error) {

	var member entities.FamilyMember

	err := r.db.
		Preload("User").
		Preload("Family").
		First(&member, entities.FamilyMember{
			UserID:   userID,
			FamilyID: familyID,
		}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// CreateMember создает новый объект члена семьи в БД
func (r *MemberRepo) CreateMember(member *entities.FamilyMember) error {
	return r.db.Create(member).Error
}

// UpdateMember обновляет объект члена семьи в БД
func (r *MemberRepo) UpdateMember(member *entities.FamilyMember) error {
	return r.db.Save(member).Error
}

// DeleteMember удаляет запись об определенном члене семьи
func (r *MemberRepo) DeleteMember(member *entities.FamilyMember) error {
	return r.db.Delete(member).Error
}
