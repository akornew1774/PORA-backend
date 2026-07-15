// ports - пакет, содержащий все порты (интерфейсы)
package ports

import (
	"pora/internal/entities"
	"time"

	"github.com/google/uuid"
)

// UserRepository содержит порты для взаимодействия с User в БД
type UserRepository interface {
	// FindByID находит пользователя по его ID
	FindByID(userID uuid.UUID) (*entities.User, error)

	// FindWithLists находит пользователя на ID вместе с его списками продуктов
	FindWithLists(userID uuid.UUID) (*entities.User, error)

	// FindWithFamilies находит пользователя на ID вместе с его семьями
	FindWithFamilies(userID uuid.UUID) (*entities.User, error)

	// FindByPhone находит пользователя по номеру телефона
	FindByPhone(phone string) (*entities.User, error)

	// FindByEmail находит пользователя по электронной почте
	FindByEmail(email string) (*entities.User, error)

	// CreateUser создает новый объект пользователя в БД
	CreateUser(user *entities.User) error

	// UpdateUser изменяет данные определенного пользователя
	UpdateUser(user *entities.User) error
}

// OTPRepository содержит порты для взаимодействия с Otp в БД
type OtpRepository interface {
	// FindByPhone находит Otp по номеру телефона
	FindByPhone(phone string) (*entities.Otp, error)

	// FindByEmail находит Otp по электронной почте
	FindByEmail(email string) (*entities.Otp, error)

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

// FamilyRepository содержит порты для
// взаимодействия с семьями в БД
type FamilyRepository interface {
	// FindByID находит семью по ID без связанный с ней сущностей
	FindByID(familyID uuid.UUID) (*entities.Family, error)

	// FindWithMembers находит семью по ID вместе с её участниками
	FindWithMembers(familyID uuid.UUID) (*entities.Family, error)

	// FindWithLists находит семью по ID вместе с её списками продуктов
	FindWithLists(familyID uuid.UUID) (*entities.Family, error)

	// FindWithAll находит семью по ID вместе со всеми связанными сущностями
	FindWithAll(familyID uuid.UUID) (*entities.Family, error)

	// FindByCode находит семью по уникальному коду для приглашения новых членов
	FindByCode(familyCode string) (*entities.Family, error)

	// CreateFamily создает новый объект семьи в БД
	CreateFamily(family *entities.Family) error

	// UpdateFamily создает новый объект семьи в БД
	UpdateFamily(family *entities.Family) error

	// DeleteFamily создает новый объект семьи в БД
	DeleteFamily(family *entities.Family) error
}

// MemberRepository содержит порты для
// взаимодействия с участниками семей в БД
type MemberRepository interface {
	// FindByUserAndFamily находит участника по его ID и ID семьи
	FindByUserAndFamily(userID uuid.UUID, familyID uuid.UUID) (
		*entities.FamilyMember, error)

	// CreateMember создает новый объект члена семьи в БД
	CreateMember(member *entities.FamilyMember) error

	// UpdateMember обновляет объект члена семьи в БД
	UpdateMember(member *entities.FamilyMember) error

	// DeleteMember удаляет запись об определенном члене семьи
	DeleteMember(member *entities.FamilyMember) error
}

// ListRepository содержит порты для
// взаимодействия со списками продуктов в БД
type ListRepository interface {
	// FindByID находит спискок покупок вместе со всеми товарами
	FindByID(listID uuid.UUID) (*entities.List, error)

	// CreateList создает новый объект списка продуктов в БД
	CreateList(list *entities.List) error

	// UpdateList обновляет существующий объект списка продуктов в БД
	UpdateList(list *entities.List) error

	// DeleteList удаляет существующий объект списка продуктов в БД
	DeleteList(list *entities.List) error
}

// ItemRepository содержит порты для
// взаимодействия с товарами в БД
type ItemRepository interface {
	// FindByID находит товар из списка продуктов по ID
	FindByID(itemID uuid.UUID) (*entities.Item, error)

	// FindItemsToRemind находит все товары, о которых
	// нужно уведомить в определенный момент времени
	FindItemsToRemind(now time.Time) ([]entities.Item, error)

	// FindItemsToUncheck находит все товары, для
	// которых нужно изменить поле checked на false
	FindItemsToUncheck(intervalAgo time.Time) ([]entities.Item, error)

	// CreateItem создает новый объект товара в БД
	CreateItem(item *entities.Item) error

	// UpdateItem обновляет существующий объект товара в БД
	UpdateItem(item *entities.Item) error

	// DeleteItem удаляет существующий объект товара в БД
	DeleteItem(item *entities.Item) error
}

// DeviceRepository содержит порты для
// взаимодействия с устройствами пользователей в БД
type DeviceRepository interface {
	// FindByUserID находит устройство пользователя по ID пользователя
	FindByUserID(userID uuid.UUID) (*entities.Device, error)

	// CreateDevice создает новый объект устройства в БД
	CreateDevice(device *entities.Device) error

	// UpdateDevice обновляет существующий объект устройства в БД
	UpdateDevice(device *entities.Device) error

	// DeleteDevice удаляет существующий объект устройства в БД
	DeleteDevice(device *entities.Device) error
}
