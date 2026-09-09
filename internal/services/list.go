// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"bytes"
	"encoding/json"
	"os"
	"pora/internal/config"
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
	"pora/internal/entities"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"
	"pora/internal/websocket"
	"sort"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// ListService - объект, содержащий методы
// для работы со списками продуктов
type ListService struct {
	userRepo   ports.UserRepository
	memberRepo ports.MemberRepository
	familyRepo ports.FamilyRepository
	listRepo   ports.ListRepository
	itemRepo   ports.ItemRepository
	hub        ports.Hub
	itemConfig config.ItemConfig
}

// NewListService создает и возвращает новый объект ListService
func NewListService(userRepo ports.UserRepository,
	memberRepo ports.MemberRepository,
	familyRepo ports.FamilyRepository,
	listRepo ports.ListRepository,
	itemRepo ports.ItemRepository,
	hub *websocket.Hub,
	itemConfig config.ItemConfig) ports.ListService {
	return &ListService{
		userRepo:   userRepo,
		memberRepo: memberRepo,
		familyRepo: familyRepo,
		listRepo:   listRepo,
		itemRepo:   itemRepo,
		hub:        hub,
		itemConfig: itemConfig,
	}
}

// CreateList создает новый список продуктов
// для пользователя/семьи и возвращает его ID
func (s *ListService) CreateList(userID uuid.UUID,
	familyID *uuid.UUID, name string) (responses.IDResponse, error) {

	var response responses.IDResponse

	listID := uuid.New()

	newList := &entities.List{
		ID:   listID,
		Name: name,
	}

	if familyID == nil || *familyID == uuid.Nil {

		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			logger.Log.Error("Ошибка при поиске пользователя: ", err)
			return response, err
		}

		if user == nil {
			logger.Log.Warn("Указанный пользователь не найден")
			return response, errors.ErrorUserNotFound
		}

		newList.UserID = &userID

	} else {

		family, err := s.familyRepo.FindWithMembers(*familyID)
		if err != nil {
			logger.Log.Error("Ошибка при поиске семьи: ", err)
			return response, err
		}

		if family == nil {
			logger.Log.Warn("Указанная семья не найдена")
			return response, errors.ErrorFamilyNotFound
		}

		newList.FamilyID = familyID

		err = s.SendChangesToFamily(family, nil, nil)
		if err != nil {
			logger.Log.Warn("Ошибка при отправке изменений членам семьи: ", err)
		}
	}

	err := s.listRepo.CreateList(newList)
	if err != nil {
		logger.Log.Error("Ошибка при сохранении списка продуктов: ", err)
		return response, err
	}

	response.ID = listID
	return response, nil
}

// GetListInfo получает полную информацию о конкретном списке продуктов
func (s *ListService) GetListInfo(listID uuid.UUID) (
	responses.ListInfo, error) {

	var response responses.ListInfo

	list, err := s.listRepo.FindByID(listID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске списка по ID")
		return response, err
	}

	if list == nil {
		logger.Log.Warn("Список с указанным ID не найден")
		return response, errors.ErrorListNotFound
	}

	sections, err := s.GetAllSections(list)

	response = responses.ListInfo{
		ID:        listID,
		Name:      list.Name,
		Sections:  sections,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}

	return response, nil
}

// AddItem добавляет один новый товар в указанный список покупок
func (s *ListService) AddItem(userID uuid.UUID, listID uuid.UUID,
	req requests.ChangeItemRequest) (responses.IDResponse, error) {

	var response responses.IDResponse

	list, err := s.listRepo.FindByID(listID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске списка по ID: ", err)
		return response, err
	}

	if list == nil {
		logger.Log.Warn("Указанный список не найден")
		return response, errors.ErrorListNotFound
	}

	itemID, err := s.createNewItem(userID, listID, req)
	if err != nil {
		return response, err
	}

	list.UpdatedAt = time.Now().UTC()

	err = s.listRepo.UpdateList(list)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении списка: ", err)
		return response, err
	}

	response.ID = itemID

	if list.FamilyID != nil && *list.FamilyID != uuid.Nil {

		family, err := s.familyRepo.FindWithMembers(*list.FamilyID)
		if err != nil {
			logger.Log.Error("Ошибка при поиске семьи: ", err)
			return response, nil
		}

		if family == nil {
			logger.Log.Warn("Указанная семья не найдена")
			return response, nil
		}

		err = s.SendChangesToFamily(family, &listID, nil)
		if err != nil {
			logger.Log.Warn("Ошибка при отправке изменений членам семьи: ", err)
		}
	}

	return response, nil
}

// AddItems добавляет несколько новых товаров в список покупок
func (s *ListService) AddItems(userID uuid.UUID,
	listID uuid.UUID, req requests.AddItemsRequest) error {

	list, err := s.listRepo.FindByID(listID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске списка по ID: ", err)
		return err
	}

	if list == nil {
		logger.Log.Warn("Указанный список не найден")
		return errors.ErrorListNotFound
	}

	for _, changeItemRequest := range req.Items {

		_, err := s.createNewItem(userID, listID, changeItemRequest)
		if err != nil {
			return err
		}
	}

	list.UpdatedAt = time.Now().UTC()

	err = s.listRepo.UpdateList(list)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении списка: ", err)
		return err
	}

	if list.FamilyID != nil && *list.FamilyID != uuid.Nil {

		family, err := s.familyRepo.FindWithMembers(*list.FamilyID)
		if err != nil {
			logger.Log.Error("Ошибка при поиске семьи: ", err)
			return nil
		}

		if family == nil {
			logger.Log.Warn("Указанная семья не найдена")
			return nil
		}

		err = s.SendChangesToFamily(family, &listID, nil)
		if err != nil {
			logger.Log.Warn("Ошибка при отправке изменений членам семьи: ", err)
		}
	}

	return nil
}

// createNewItem создает новую сущность товара со всеми полями и сохраняет ее в БД
func (s *ListService) createNewItem(userID uuid.UUID, listID uuid.UUID,
	req requests.ChangeItemRequest) (uuid.UUID, error) {

	itemID := uuid.New()

	item := &entities.Item{
		ID:     itemID,
		ListID: listID,

		Name:    req.Name,
		Section: req.Section,

		Quantity: req.Quantity,
		Unit:     req.Unit,

		Priority: req.Priority,
		Urgent:   req.Urgent,

		Checked:         req.Checked,
		RemindEveryDays: req.RemindEveryDays,

		AddedByID: userID,
	}

	if item.Checked {
		checkedAt := time.Now().UTC()
		item.CheckedAt = &checkedAt
		item.TimesBought++
	}

	if item.RemindEveryDays != nil && *item.RemindEveryDays <= 0 {
		logger.Log.Warn("Некорректный формат RemindEveryDays")
		return itemID, errors.ErrorInvalidInput
	}

	if item.RemindEveryDays != nil {

		t, err := time.Parse("15:04", s.itemConfig.DefaultReminderTime)
		if err != nil {
			logger.Log.Error("Ошибка при парсинге времени: ", err)
			return itemID, err
		}

		nextDate := time.Now().UTC().AddDate(0, 0, *item.RemindEveryDays)

		nextReminder := time.Date(
			nextDate.Year(),
			nextDate.Month(),
			nextDate.Day(),
			t.Hour(),
			t.Minute(),
			0,
			0,
			nextDate.Location(),
		)

		item.NextReminderAt = &nextReminder
	}

	err := s.itemRepo.CreateItem(item)
	if err != nil {
		logger.Log.Error("Ошибка при создании товара в БД: ", err)
		return itemID, err
	}

	return itemID, nil
}

// DeleteList удаляет один конкретный список продуктов
func (s *ListService) DeleteList(listID uuid.UUID) error {

	list, err := s.listRepo.FindByID(listID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске списка продуктов: ", err)
		return err
	}

	if list == nil {
		logger.Log.Warn("Указанный список продуктов не найден")
		return errors.ErrorListNotFound
	}

	err = s.listRepo.DeleteList(list)
	if err != nil {
		logger.Log.Error("Ошибка при удалении списка продуктов: ", err)
		return err
	}

	if list.FamilyID != nil && *list.FamilyID != uuid.Nil {

		family, err := s.familyRepo.FindWithMembers(*list.FamilyID)
		if err != nil {
			logger.Log.Error("Ошибка при поиске семьи: ", err)
			return nil
		}

		if family == nil {
			logger.Log.Warn("Указанная семья не найдена")
			return nil
		}

		err = s.SendChangesToFamily(family, nil, nil)
		if err != nil {
			logger.Log.Warn("Ошибка при отправке изменений членам семьи: ", err)
		}
	}

	return nil
}

// GetAllSections делит все товары на секции и
// переводит их в нужный для response формат
func (s *ListService) GetAllSections(list *entities.List) (
	[]responses.SectionInfo, error) {

	sections := make(map[string][]responses.ItemInfo)

	for _, item := range list.Items {

		itemInfo, err := s.СonvertToItemInfo(&item, list.FamilyID)
		if err != nil {
			logger.Log.Warn("Ошибка при получении информации о товаре")
			continue
		}

		if item.Section != "" {

			sections[item.Section] = append(
				sections[item.Section],
				itemInfo,
			)

		} else {

			sections[entities.OthersSectionName] = append(
				sections[entities.OthersSectionName],
				itemInfo,
			)
		}
	}

	names := make([]string, 0, len(sections))
	collator := collate.New(language.Russian)

	for name := range sections {
		names = append(names, name)
	}

	sort.Slice(names, func(i, j int) bool {
		if names[i] == "" {
			return false
		}
		if names[j] == "" {
			return true
		}
		return collator.CompareString(names[i], names[j]) < 0
	})

	result := make([]responses.SectionInfo, 0, len(names))

	for _, name := range names {

		sectionItems, err := s.sortListItems(sections[name], collator)
		if err != nil {
			logger.Log.Error("Ошибка при сортировке модели товара: ", err)
			return result, err
		}

		result = append(result, responses.SectionInfo{
			Name:  name,
			Items: sectionItems,
		})
	}

	return result, nil
}

// GetHighestPrioritySections отбирает товары высшего приоритета, делит
// их на секции и приводит к нужному для response формату
func (s *ListService) GetHighestPrioritySections(
	list *entities.List) ([]responses.SectionInfo, error) {

	sections := make(map[string][]responses.ItemInfo)

	for _, item := range list.Items {

		if item.Priority != entities.HighestItemPriority && !item.Urgent {
			continue
		}

		itemInfo, err := s.СonvertToItemInfo(&item, list.FamilyID)
		if err != nil {
			logger.Log.Warn("Ошибка при получении информации о товаре")
			continue
		}

		if item.Section != "" {

			sections[item.Section] = append(
				sections[item.Section],
				itemInfo,
			)

		} else {

			sections[entities.OthersSectionName] = append(
				sections[entities.OthersSectionName],
				itemInfo,
			)
		}
	}

	names := make([]string, 0, len(sections))
	collator := collate.New(language.Russian)

	for name := range sections {
		names = append(names, name)
	}

	sort.Slice(names, func(i, j int) bool {
		if names[i] == "" {
			return false
		}
		if names[j] == "" {
			return true
		}
		return collator.CompareString(names[i], names[j]) < 0
	})

	result := make([]responses.SectionInfo, 0, len(names))

	for _, name := range names {

		sectionItems, err := s.sortListItems(sections[name], collator)
		if err != nil {
			logger.Log.Error("Ошибка при сортировке модели товара: ", err)
			return result, err
		}

		result = append(result, responses.SectionInfo{
			Name:  name,
			Items: sectionItems,
		})
	}

	return result, nil
}

// sortListItems сортирует товары списка
// (сначала по Proirity, потом по алфавиту)
func (s *ListService) sortListItems(items []responses.ItemInfo,
	collator *collate.Collator) ([]responses.ItemInfo, error) {

	sort.SliceStable(items, func(i, j int) bool {
		return bytes.Compare(items[i].ID[:], items[j].ID[:]) < 0
	})

	sort.SliceStable(items, func(i, j int) bool {
		return collator.CompareString(
			items[i].Name,
			items[j].Name,
		) < 0
	})

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Priority > items[j].Priority
	})

	return items, nil
}

// СonvertToItemInfo переводит сущности Item из типа entities.Item в формат
// responses.ItemInfo, при этом добавляя информацию о добавившем его пользователе
func (s *ListService) СonvertToItemInfo(item *entities.Item,
	familyID *uuid.UUID) (responses.ItemInfo, error) {

	itemInfo := responses.ItemInfo{
		ID:              item.ID,
		Name:            item.Name,
		Quantity:        item.Quantity,
		Unit:            item.Unit,
		Priority:        item.Priority,
		Urgent:          item.Urgent,
		Checked:         item.Checked,
		RemindEveryDays: item.RemindEveryDays,
	}

	if item.AddedByID == uuid.Nil || familyID == nil || *familyID == uuid.Nil {
		return itemInfo, nil
	}

	member, err := s.memberRepo.
		FindByUserAndFamily(item.AddedByID, *familyID)

	if err != nil {
		logger.Log.Error("Ошибка при поиске члена семьи: ", err)
		return itemInfo, err
	}

	if member == nil || member.User == nil {
		logger.Log.Warn("Указанный член семьи не найден")
		return itemInfo, nil
	}

	user := member.User

	memberInfo := &responses.FamilyMemberInfo{

		UserInfo: responses.UserInfo{
			ID:       user.ID,
			Name:     user.Name,
			Surname:  user.Surname,
			ImageURL: os.Getenv("BASE_URL") + user.ImageURL,
		},

		JoinedAt: member.JoinedAt,
		Color:    string(member.Color),
	}

	itemInfo.AddedBy = memberInfo

	return itemInfo, nil
}

// SendChangesToFamily отправляет уведомление об изменении
// семьи/списка/товара через Websocket
func (s *ListService) SendChangesToFamily(family *entities.Family,
	listID *uuid.UUID, itemID *uuid.UUID) error {

	message := websocket.SendChangesMessage{
		FamilyID: &family.ID,
		ListID:   listID,
		ItemID:   itemID,
	}

	data, err := json.Marshal(message)
	if err != nil {
		logger.Log.Error("Ошибка при сериализации сообщения в JSON: ", err)
		return err
	}

	for _, member := range family.Members {
		s.hub.SendToUser(member.UserID, data)
	}

	return nil
}
