// ports - пакет, содержащий все порты (интерфейсы)
package ports

import (
	"context"
	"io"
	"pora/internal/dto/notifications"
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
	"pora/internal/entities"

	"github.com/google/uuid"
)

// UserService содержит порты для методов, необходимых для
// работы с пользователем и его профилем
type UserService interface {
	// UpdateUserInfo обновляет (или дополняет) информацию
	// о каком-то конкретном пользователе
	UpdateUserInfo(userID uuid.UUID,
		req *requests.UpdateUserRequest) error

	// UpdateDevice обновляет токен устройства у пользователя
	UpdateDevice(userID uuid.UUID,
		deviceToken string, deviceType string) error

	// SaveImage вызывает FileService для сохранения изображения в
	// хранилище, сохраняет ссылку на него в профиле пользователя в БД
	SaveImage(userID uuid.UUID,
		file io.Reader, fileName string) (string, error)

	// GetMyInfo получает информацию о текущем пользователе
	GetMyInfo(userID uuid.UUID) (responses.GetMyInfoResponse, error)
}

// OtpService содержит порты для методов отправки и подтверждения OTP
type OtpService interface {
	// SendOtp отправляет Otp-код по номеру телефона / почте и сохраняет его
	SendOtp(rawPhone string, email string) error

	// VerifyOtp сравнивает полученный Otp-код c сохраненным в БД
	VerifyOtp(req requests.VerifyOtpRequest) (
		responses.VerifyOtpResponse, error)
}

// AuthService содержит порты для методов, необходимых для авторизации
type AuthService interface {
	// CheckUser находит пользователя по телефону / email и проверяет его
	// статус (notFound / notRegistered / registered)
	CheckUser(rawPhone string, email string) (responses.IsUserResponse, error)

	// GetNewTokens получает новые access и refresh токены для пользователя
	GetNewTokens(refreshTokenStr string) (responses.RefreshResponse, error)

	// RemoveRefreshToken удаляет refresh-токен по ID пользователя
	RemoveRefreshToken(userID uuid.UUID) error
}

// TokenService содержит порты для методов создания/обновления токенов
type TokenService interface {
	// CreateTokens создает новый access-токен,
	// обновляет (или создает новый) refresh-токен,
	// после чего возвращает access и refresh токены.
	CreateTokens(user *entities.User) (string, string, error)

	// DecodeAccessToken декодирует полученный access-токен,
	// проверяет его валидность
	DecodeAccessToken(rawToken string) (uuid.UUID, error)
}

// FileService содержит порты для методов для
// работы с изображениями в хранилище
type FileService interface {
	// SaveProfileImage сохраняет изображение профиля пользователя
	// в хранилище и возвращает ссылку на него
	SaveProfileImage(file io.Reader,
		fileName string, userID uuid.UUID) (string, error)
}

// FamilyService содержит порты для методов для
// создания семей и просмотра информации о них
type FamilyService interface {
	// GetFamilies получает информацию в всех семьях текущего пользователя
	GetFamilies(userID uuid.UUID) (responses.GetFamiliesResponse, error)

	// GetLists получает информацию о списках продуктов конкретной семьи
	GetLists(familyID uuid.UUID) (responses.GetListsResponse, error)

	// CreateFamily создает новую семью и возвращает её ID
	CreateFamily(userID uuid.UUID, name string) (responses.IDResponse, error)

	// AddMember добавляет к существующей семье еще одного участника
	AddMember(userID uuid.UUID, familyCode string) error

	// GetFamilyLink получает Link-код семьи
	// вместе с сылкой для вступления в неё
	GetFamilyLink(familyID uuid.UUID) (responses.GetFamilyLinkResponse, error)

	// DeleteFamily удаляет одну конкретную семью
	DeleteFamily(familyID uuid.UUID) error
}

// ListService содержит порты для методов для
// просмотра и изменения списков покупок
type ListService interface {
	// CreateList создает новый список продуктов
	// для пользователя/семьи и возвращает его ID
	CreateList(userID uuid.UUID, familyID *uuid.UUID,
		name string) (responses.IDResponse, error)

	// GetListInfo получает полную информацию о конкретном списке продуктов
	GetListInfo(listID uuid.UUID) (responses.ListInfo, error)

	// AddItem добавляет один новый товар в указанный список покупок
	AddItem(userID uuid.UUID, listID uuid.UUID,
		req requests.ChangeItemRequest) (responses.IDResponse, error)

	// DeleteList удаляет один конкретный список продуктов
	DeleteList(listID uuid.UUID) error

	// GetAllSections делит все товары на секции и
	// переводит их в нужный для response формат
	GetAllSections(list *entities.List) (
		[]responses.SectionInfo, error)

	// GetHighestPrioritySections отбирает товары высшего приоритета, делит
	// их на секции и приводит к нужному для response формату
	GetHighestPrioritySections(list *entities.List) (
		[]responses.SectionInfo, error)

	// СonvertToItemInfo переводит сущности Item из типа entities.Item в формат
	// responses.ItemInfo, при этом добавляя информацию о добавившем его пользователе
	СonvertToItemInfo(item *entities.Item,
		familyID *uuid.UUID) (responses.ItemInfo, error)

	// SendChangesToFamily отправляет уведомление об изменении
	// семьи/списка/товара через Websocket
	SendChangesToFamily(family *entities.Family,
		listID *uuid.UUID, itemID *uuid.UUID) error
}

// ItemService содержит порты для методов для
// изменения и удаления продуктов, а также прочих
// взаимодействий с данной сущностью
type ItemService interface {
	// GetItemInfo получает информацию о конкретном товаре
	GetItemInfo(itemID uuid.UUID) (responses.ItemInfo, error)

	// ChangeItem изменяет поля определенного товара
	ChangeItem(itemID uuid.UUID, req requests.ChangeItemRequest) error

	// DeleteItem удаляет один товар из списка продуктов
	DeleteItem(itemID uuid.UUID) error

	// MarkAsBought ставит значение true в поле Checked у товара
	MarkAsBought(itemID uuid.UUID, checked bool) error

	// NotifyMembers уведомляет указанных членов семьи об
	// определенном продукте. К уведомлению можно прикрепить сообщение
	NotifyMembers(ctx context.Context, userID uuid.UUID,
		itemID uuid.UUID, req requests.NotifyMembersRequest) error

	// ProcessReminders находит товары, о которых нужно уведомить
	// пользователя, и вызывает методы для отправки уведомлений
	ProcessReminders(ctx context.Context) error

	// ProcessCheckedItems находит товары, для которых
	// необходимо убрать значение true в поле Checked /
	// удалить товар (если нет значения RemindEveryDays)
	ProcessCheckedItems(ctx context.Context) error
}

// PushService содержит порты для методов для
// отправки Push-уведомлений пользователям
type PushService interface {
	// SendToUser отправляет Push-уведомление указанному пользователю
	SendToUser(ctx context.Context, userID uuid.UUID,
		notification notifications.Notification) error

	// SendItemNotification отправляет пользователю
	// Push-уведомление о конкретном товаре от члена семьи
	SendItemNotification(ctx context.Context,
		itemNotification notifications.ItemNotification) error

	// SendItemReminder отправляет пользователю
	// Push-уведомление с напоминанием о конкретном товаре
	SendItemReminder(ctx context.Context,
		itemReminder notifications.ItemReminder) error
}

// StatisticsService содержит порты для методов
// для предоставления пользователям статистики
type StatisticsService interface {
	// GetUserProducts получает все созданные пользователем товары
	GetUserProducts(userID uuid.UUID) (
		responses.GetUserProductsResponse, error)

	// GetBrief получает все быстро кончающиеся товары пользователя
	GetBrief(userID uuid.UUID) (
		responses.GetBriefResponse, error)

	// SaveBrief сохраняет быстро кончающиеся товары пользователя
	SaveBrief(userID uuid.UUID,
		req requests.SaveBriefRequest) error

	// GetLoginTimes получает все время входа пользователя в приложение
	GetLoginTimes(userID uuid.UUID) (
		responses.LoginTimesResponse, error)
}
