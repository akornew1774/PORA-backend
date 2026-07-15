// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"context"
	"pora/internal/config"
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"
	"slices"
	"time"

	"github.com/google/uuid"
)

// ItemService - объект, содержащий методы для
// работы с товарами в списках покупок
type ItemService struct {
	familyRepo ports.FamilyRepository
	listRepo   ports.ListRepository
	itemRepo   ports.ItemRepository

	listService ports.ListService
	pushService ports.PushService

	itemConfig config.ItemConfig
}

// NewItemService создает и возвращает новый объект ItemService
func NewItemService(familyRepo ports.FamilyRepository,
	listRepo ports.ListRepository,
	itemRepo ports.ItemRepository,
	listService ports.ListService,
	pushService ports.PushService,
	itemConfig config.ItemConfig) ports.ItemService {
	return &ItemService{
		familyRepo:  familyRepo,
		listRepo:    listRepo,
		itemRepo:    itemRepo,
		listService: listService,
		pushService: pushService,
		itemConfig:  itemConfig,
	}
}

// GetItemInfo получает информацию о конкретном товаре
func (s *ItemService) GetItemInfo(itemID uuid.UUID) (
	responses.ItemInfo, error) {

	var response responses.ItemInfo

	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске товара по ID: ", err)
		return response, err
	}

	if item == nil {
		logger.Log.Warn("Указанный товар не найден")
		return response, errors.ErrorItemNotFound
	}

	list, err := s.listRepo.FindByID(item.ListID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске списка по ID: ", err)
		return response, err
	}

	if list == nil {
		logger.Log.Warn("Указанный список не найден")
		return response, errors.ErrorListNotFound
	}

	response, err = s.listService.СonvertToItemInfo(item, list.FamilyID)
	if err != nil {
		logger.Log.Warn("Ошибка при конвертировании товара в нужный формат: ", err)
		return response, err
	}

	return response, nil
}

// ChangeItem изменяет поля определенного товара
func (s *ItemService) ChangeItem(itemID uuid.UUID,
	req requests.ChangeItemRequest) error {

	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске товара по ID: ", err)
		return err
	}

	if item == nil {
		logger.Log.Warn("Указанный товар не найден")
		return errors.ErrorItemNotFound
	}

	item.Name = req.Name
	item.Section = req.Section

	item.Quantity = req.Quantity
	item.Unit = req.Unit

	item.Priority = req.Priority
	item.Urgent = req.Urgent

	item.Checked = req.Checked
	item.RemindEveryDays = req.RemindEveryDays

	if item.Checked == true {
		checkedAt := time.Now().UTC()
		item.CheckedAt = &checkedAt
	}

	if item.RemindEveryDays != nil && *item.RemindEveryDays <= 0 {
		logger.Log.Warn("Некорректный формат RemindEveryDays")
		return errors.ErrorInvalidInput
	}

	if item.RemindEveryDays != nil {

		t, err := time.Parse("15:04", s.itemConfig.DefaultReminderTime)
		if err != nil {
			logger.Log.Error("Ошибка при парсинге времени: ", err)
			return err
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

	item.UpdatedAt = time.Now().UTC()

	err = s.itemRepo.UpdateItem(item)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении товара: ", err)
		return err
	}

	list, err := s.listRepo.FindByID(item.ListID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске списка по ID: ", err)
		return err
	}

	if list == nil {
		logger.Log.Warn("Указанный список не найден")
		return errors.ErrorListNotFound
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

		err = s.listService.SendChangesToFamily(family, &list.ID, &itemID)
		if err != nil {
			logger.Log.Warn("Ошибка при отпраке изменений семье: ", err)
		}
	}

	return nil
}

// DeleteItem удаляет один товар из списка продуктов
func (s *ItemService) DeleteItem(itemID uuid.UUID) error {

	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске товара по ID: ", err)
		return err
	}

	if item == nil {
		logger.Log.Warn("Товар с указанный ID не найден: ")
		return nil
	}

	err = s.itemRepo.DeleteItem(item)
	if err != nil {
		logger.Log.Error("Ошибюка при удалении товара: ", err)
		return err
	}

	list, err := s.listRepo.FindByID(item.ListID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске списка по ID: ", err)
		return err
	}

	if list == nil {
		logger.Log.Warn("Указанный список не найден")
		return errors.ErrorListNotFound
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

		err = s.listService.SendChangesToFamily(family, &list.ID, nil)
		if err != nil {
			logger.Log.Warn("Ошибка при отпраке изменений семье: ", err)
		}
	}

	return nil
}

// MarkAsBought ставит значение true в поле Checked у товара
func (s *ItemService) MarkAsBought(itemID uuid.UUID) error {

	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске товара по ID: ", err)
		return err
	}

	if item == nil {
		logger.Log.Warn("Указанный товар не найден")
		return errors.ErrorItemNotFound
	}

	item.Checked = true
	item.UpdatedAt = time.Now().UTC()

	checkedAt := time.Now().UTC()
	item.CheckedAt = &checkedAt

	err = s.itemRepo.UpdateItem(item)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении товара: ", err)
		return err
	}

	list, err := s.listRepo.FindByID(item.ListID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске списка по ID: ", err)
		return err
	}

	if list == nil {
		logger.Log.Warn("Указанный список не найден")
		return errors.ErrorListNotFound
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

		err = s.listService.SendChangesToFamily(family, &list.ID, &itemID)
		if err != nil {
			logger.Log.Warn("Ошибка при отпраке изменений семье: ", err)
		}
	}

	return nil
}

// NotifyMembers уведомляет указанных членов семьи об
// определенном продукте. К уведомлению можно прикрепить сообщение
func (s *ItemService) NotifyMembers(ctx context.Context, userID uuid.UUID,
	itemID uuid.UUID, req requests.NotifyMembersRequest) error {

	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске товара по ID: ", err)
		return err
	}

	if item == nil {
		logger.Log.Warn("Указанный товар не найден")
		return errors.ErrorItemNotFound
	}

	list, err := s.listRepo.FindByID(item.ListID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске списка по ID: ", err)
		return err
	}

	if list == nil {
		logger.Log.Warn("Указанный список не найден")
		return errors.ErrorListNotFound
	}

	if list.FamilyID == nil || *list.FamilyID == uuid.Nil {
		logger.Log.Warn("Список не принадлежит ни одной семье")
		return errors.ErrorForbidden
	}

	family, err := s.familyRepo.FindWithMembers(*list.FamilyID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске семьи: ", err)
		return err
	}

	if family == nil {
		logger.Log.Warn("Указанная семья не найдена")
		return errors.ErrorFamilyNotFound
	}

	for _, member := range family.Members {

		user := member.User

		if user == nil {
			logger.Log.Error("Член семьи не найден")
			continue
		}

		if user.ID == userID {
			continue
		}

		if len(req.To) == 0 || slices.Contains(req.To, user.ID) {

			err := s.pushService.SendItemNotification(
				ctx,
				userID,
				itemID,
				list.ID,
				family.ID,
				req.Message,
			)
			if err != nil {
				logger.Log.Warn("Ошибка при отправке Push-уведомления", err)
				continue
			}
		}
	}

	return nil
}

// ProcessReminders находит товары, о которых нужно уведомить
// пользователя, и вызывает методы для отправки уведомлений
func (s *ItemService) ProcessReminders(ctx context.Context) error {

	items, err := s.itemRepo.FindItemsToRemind(time.Now().UTC())
	if err != nil {
		logger.Log.Error("Ошибка при поиске товаров для уведомления: ", err)
		return err
	}

	for _, item := range items {

		list, err := s.listRepo.FindByID(item.ListID)
		if err != nil {
			logger.Log.Error("Ошибка при поиске списка по ID: ", err)
			return err
		}

		if list == nil {
			logger.Log.Warn("Указанный список не найден")
			return errors.ErrorListNotFound
		}

		s.remindFamilyMembers(
			ctx,
			list.UserID,
			list.FamilyID,
			item.ListID,
			item.ID,
		)

		nextReminder := item.NextReminderAt.AddDate(0, 0, *item.RemindEveryDays)

		for !nextReminder.After(time.Now().UTC()) {
			nextReminder = nextReminder.AddDate(0, 0, *item.RemindEveryDays)
		}

		item.NextReminderAt = &nextReminder

		err = s.itemRepo.UpdateItem(&item)
		if err != nil {
			logger.Log.Error("Ошибка при обновлении товара: ", err)
		}
	}

	logger.Log.Info("Напоминание о товарах выполнено")
	return nil
}

// remindFamilyMembers вызывает метод PushService для
// отправки напоминаний о товаре всем членам семьи
func (s *ItemService) remindFamilyMembers(ctx context.Context, userID *uuid.UUID,
	familyID *uuid.UUID, listID uuid.UUID, itemID uuid.UUID) {

	if familyID == nil || *familyID == uuid.Nil {

		err := s.pushService.SendItemReminder(
			ctx,
			*userID,
			itemID,
			listID,
			nil,
		)

		if err != nil {
			logger.Log.Warn("Ошибка при отправке Push-уведомления: ", err)
		}
	}

	family, err := s.familyRepo.FindWithMembers(*familyID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске семьи: ", err)
	}

	if family == nil {
		logger.Log.Warn("Указанная семья не найдена")
	}

	for _, member := range family.Members {

		err := s.pushService.SendItemReminder(
			ctx,
			member.UserID,
			itemID,
			listID,
			familyID,
		)

		if err != nil {
			logger.Log.Warn("Ошибка при отправке Push-уведомления: ", err)
		}
	}
}

// ProcessCheckedItems находит товары, для которых
// необходимо убрать значение true в поле Checked /
// удалить товар (если нет значения RemindEveryDays)
func (s *ItemService) ProcessCheckedItems(ctx context.Context) error {

	intervalHours := s.itemConfig.UncheckIntervalHours
	intervalAgo := time.Now().UTC().Add(-time.Duration(intervalHours) * time.Hour)

	items, err := s.itemRepo.FindItemsToUncheck(intervalAgo)
	if err != nil {
		logger.Log.Error("Ошибка при поиске товаров для измненения поля Checked: ", err)
		return err
	}

	for _, item := range items {

		if item.RemindEveryDays == nil {

			err := s.itemRepo.DeleteItem(&item)
			if err != nil {
				logger.Log.Error("Ошибка при удалении товара: ", err)
				return err
			}

			continue
		}

		item.Checked = false

		err := s.itemRepo.UpdateItem(&item)
		if err != nil {
			logger.Log.Error("Ошибка при удалении товара: ", err)
			return err
		}
	}

	logger.Log.Info("Обработка купленных товаров выполнена")
	return nil
}
