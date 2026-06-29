// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"pora/internal/dto/responses"
	"pora/internal/entities"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"
	"sort"

	"github.com/google/uuid"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// ListService - объект, содержащий методы
// для работы со списками продуктов
type ListService struct {
	memberRepo ports.MemberRepository
	listRepo   ports.ListRepository
	itemRepo   ports.ItemRepository
}

// NewListService создает и возвращает новый объект ListService
func NewListService(memberRepo ports.MemberRepository,
	listRepo ports.ListRepository,
	itemRepo ports.ItemRepository) ports.ListService {
	return &ListService{
		memberRepo: memberRepo,
		listRepo:   listRepo,
		itemRepo:   itemRepo,
	}
}

// GetAllSections делит все товары на секции и
// переводит их в нужный для response формат
func (s *ListService) GetAllSections(list *entities.List) (
	[]responses.SectionInfo, error) {

	sections := make(map[string][]responses.ItemInfo)

	for _, item := range list.Items {

		itemInfo, err := s.convertToItemInfo(&item, list.FamilyID)
		if err != nil {
			logger.Log.Warn("Ошибка при получении информации о товаре")
			continue
		}

		sections[item.Section] = append(sections[item.Section], itemInfo)
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

		itemInfo, err := s.convertToItemInfo(&item, list.FamilyID)
		if err != nil {
			logger.Log.Warn("Ошибка при получении информации о товаре")
			continue
		}

		sections[item.Section] = append(sections[item.Section], itemInfo)
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

// convertToItemInfo переводит сущности Item из типа entities.Item в формат
// responses.ItemInfo, при этом добавляя информацию о добавившем его пользователе
func (s *ListService) convertToItemInfo(item *entities.Item,
	familyID uuid.UUID) (responses.ItemInfo, error) {

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

	if item.AddedByID == nil || familyID == uuid.Nil {
		return itemInfo, nil
	}

	member, err := s.memberRepo.
		FindByUserAndFamily(*item.AddedByID, familyID)

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
