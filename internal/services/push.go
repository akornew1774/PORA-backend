// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"context"
	"fmt"
	"pora/internal/config"
	"pora/internal/dto/notifications"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"

	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
)

// PushService - объект, содержащий методы для
// отправки Push-уведомлений пользователям
type PushService struct {
	deviceRepo     ports.DeviceRepository
	firebaseClient ports.FirebaseClient
	cfg            config.PushConfig
}

// NewPushService создает и возвращает новый объект PushService
func NewPushService(deviceRepo ports.DeviceRepository,
	firebaseClient ports.FirebaseClient,
	cfg config.PushConfig) ports.PushService {
	return &PushService{
		deviceRepo:     deviceRepo,
		firebaseClient: firebaseClient,
		cfg:            cfg,
	}
}

// SendToUser отправляет Push-уведомление указанному пользователю
func (s *PushService) SendToUser(ctx context.Context,
	userID uuid.UUID, notification notifications.Notification) error {

	device, err := s.deviceRepo.FindByUserID(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске устройства пользователя: ", err)
		return err
	}

	if device == nil {
		logger.Log.Warn("Указанное устройство не найдено")
		return errors.ErrorUserNotFound
	}

	message := &messaging.Message{
		Token: device.DeviceToken,

		Notification: &messaging.Notification{
			Title: notification.Title,
			Body:  notification.Body,
		},

		Data: notification.Data,
	}

	_, err = s.firebaseClient.Send(ctx, message)
	if err != nil {
		logger.Log.Error("Ошибка при отправке уведомления через firebase: ", err)
		return err
	}

	return nil
}

// SendItemNotification отправляет пользователю
// Push-уведомление о конкретном товаре от члена семьи
func (s *PushService) SendItemNotification(ctx context.Context,
	itemNotification notifications.ItemNotification) error {

	var title string
	var body string

	if itemNotification.Message != "" {

		title = itemNotification.Message

		body = fmt.Sprintf(
			"%s просит купить %s, %s %s",
			itemNotification.AuthorName,
			itemNotification.ItemName,
			itemNotification.ItemQuantity,
			itemNotification.ItemUnit,
		)

	} else {

		title = fmt.Sprintf(
			"Купи %s, %s %s",
			itemNotification.ItemName,
			itemNotification.ItemQuantity,
			itemNotification.ItemUnit,
		)

		body = fmt.Sprintf(
			"%s отправил(а) сообщение о товаре из списка \"%s\"",
			itemNotification.AuthorName,
			itemNotification.ListName,
		)
	}

	data := make(map[string]string)

	data["type"] = s.cfg.NotificationType
	data["screen"] = s.cfg.NotificationScreen

	data["family-id"] = itemNotification.FamilyID.String()
	data["list-id"] = itemNotification.ListID.String()
	data["item-id"] = itemNotification.ItemID.String()

	notification := notifications.Notification{
		Title: title,
		Body:  body,
		Data:  data,
	}

	return s.SendToUser(ctx, itemNotification.UserID, notification)
}

// SendItemReminder отправляет пользователю
// Push-уведомление с напоминанием о конкретном товаре
func (s *PushService) SendItemReminder(ctx context.Context,
	itemReminder notifications.ItemReminder) error {

	title := fmt.Sprintf(
		"Не забудьте купить %s, %s %s",
		itemReminder.ItemName,
		itemReminder.ItemQuantity,
		itemReminder.ItemUnit,
	)

	body := fmt.Sprintf(
		"Напоминание о товаре из списка \"%s\"",
		itemReminder.ListName,
	)

	data := make(map[string]string)

	data["type"] = s.cfg.NotificationType
	data["screen"] = s.cfg.NotificationScreen

	data["list-id"] = itemReminder.ListID.String()
	data["item-id"] = itemReminder.ItemID.String()

	familyID := itemReminder.FamilyID

	if familyID != nil && *familyID != uuid.Nil {
		data["family-id"] = familyID.String()
	}

	notification := notifications.Notification{
		Title: title,
		Body:  body,
		Data:  data,
	}

	return s.SendToUser(ctx, itemReminder.UserID, notification)
}
