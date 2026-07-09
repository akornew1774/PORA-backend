// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
	"pora/internal/entities"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"
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
}

// NewListService создает и возвращает новый объект ListService
func NewListService(userRepo ports.UserRepository,
	memberRepo ports.MemberRepository,
	familyRepo ports.FamilyRepository,
	listRepo ports.ListRepository,
	itemRepo ports.ItemRepository) ports.ListService {
	return &ListService{
		userRepo:   userRepo,
		memberRepo: memberRepo,
		familyRepo: familyRepo,
		listRepo:   listRepo,
		itemRepo:   itemRepo,
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

		family, err := s.familyRepo.FindByID(*familyID)
		if err != nil {
			logger.Log.Error("Ошибка при поиске семьи: ", err)
			return response, err
		}

		if family == nil {
			logger.Log.Warn("Указанная семья не найдена")
			return response, errors.ErrorFamilyNotFound
		}

		newList.FamilyID = familyID
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

	list, err := s.listRepo.FindByID(listID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске списка по ID: ", err)
		return response, err
	}

	if list == nil {
		logger.Log.Warn("Указанный список не найден")
		return response, errors.ErrorListNotFound
	}

	err = s.itemRepo.CreateItem(item)
	if err != nil {
		logger.Log.Error("Ошибка при создании товара в БД: ", err)
		return response, err
	}

	list.UpdatedAt = time.Now().UTC()

	err = s.listRepo.UpdateList(list)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении списка: ", err)
		return response, err
	}

	response.ID = itemID
	return response, nil
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
		result = append(result, responses.SectionInfo{
			Name:  name,
			Items: sections[name],
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
		result = append(result, responses.SectionInfo{
			Name:  name,
			Items: sections[name],
		})
	}

	return result, nil
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
			ImageURL: user.ImageURL,
		},

		JoinedAt: member.JoinedAt,
		Color:    string(member.Color),
	}

	itemInfo.AddedBy = memberInfo

	return itemInfo, nil
}
