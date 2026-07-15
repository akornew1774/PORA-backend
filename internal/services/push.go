// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"context"
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
	userID uuid.UUID, itemID uuid.UUID, listID uuid.UUID,
	familyID uuid.UUID, message string) error {

	// TODO: сделать формирование уведомления
	notification := notifications.Notification{}

	return s.SendToUser(ctx, userID, notification)
}

// SendItemReminder отправляет пользователю
// Push-уведомление с напоминанием о конкретном товаре
func (s *PushService) SendItemReminder(ctx context.Context,
	userID uuid.UUID, itemID uuid.UUID, listID uuid.UUID,
	familyID *uuid.UUID) error {

	// TODO: сделать формирование напомнинания
	notification := notifications.Notification{}

	return s.SendToUser(ctx, userID, notification)
}
