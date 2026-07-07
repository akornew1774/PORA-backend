// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"pora/internal/dto/requests"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"
	"slices"
	"time"

	"github.com/google/uuid"
)

type ItemService struct {
	familyRepo ports.FamilyRepository
	listRepo   ports.ListRepository
	itemRepo   ports.ItemRepository
}

// NewItemService создает и возвращает новый объект ItemService
func NewItemService(familyRepo ports.FamilyRepository,
	listRepo ports.ListRepository,
	itemRepo ports.ItemRepository) ports.ItemService {
	return &ItemService{
		familyRepo: familyRepo,
		listRepo:   listRepo,
		itemRepo:   itemRepo,
	}
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

	return nil
}

// NotifyMembers уведомляет указанных членов семьи об
// определенном продукте. К уведомлению можно прикрепить сообщение
func (s *ItemService) NotifyMembers(userID uuid.UUID,
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

			// TODO: сделать уведомление членов семьи

		}
	}

	return nil
}
