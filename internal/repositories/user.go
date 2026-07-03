// repositories - пакет с методами для работы с БД
package repositories

import (
	"errors"
	"pora/internal/entities"
	"pora/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepo - объект, содержащий методы для работы с сущностью пользователя
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepository создает и возвращает новый объект UserRepo
func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &UserRepo{db: db}
}

// FindByID находит пользователя по его ID
func (r *UserRepo) FindByID(userID uuid.UUID) (*entities.User, error) {
	var user entities.User

	err := r.db.First(&user, userID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindWithLists находит пользователя на ID вместе с его списками продуктов
func (r *UserRepo) FindWithLists(userID uuid.UUID) (*entities.User, error) {
	var user entities.User

	err := r.db.
		Preload("Lists").
		First(&user, userID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindWithFamilies находит пользователя на ID вместе с его семьями
func (r *UserRepo) FindWithFamilies(userID uuid.UUID) (*entities.User, error) {
	var user entities.User

	err := r.db.
		Preload("Memberships").
		Preload("Memberships.Family").
		Preload("Memberships.Family.Members").
		Preload("Memberships.Family.Members.User").
		First(&user, userID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone находит пользователя по номеру телефона
func (r *UserRepo) FindByPhone(phone string) (*entities.User, error) {
	var user entities.User

	err := r.db.
		Where("phone = ?", phone).
		First(&user).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail находит пользователя по электронной почте
func (r *UserRepo) FindByEmail(email string) (*entities.User, error) {
	var user entities.User

	err := r.db.
		Where("email = ?", email).
		First(&user).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser создает новый объект пользователя в БД
func (r *UserRepo) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
}

// UpdateUser изменяет данные определенного пользователя
func (r *UserRepo) UpdateUser(user *entities.User) error {
	return r.db.Save(user).Error
}
