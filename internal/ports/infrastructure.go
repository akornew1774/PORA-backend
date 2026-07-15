// ports - пакет, содержащий все порты (интерфейсы)
package ports

import (
	"context"
	"io"

	"firebase.google.com/go/v4/messaging"
)

// Storage содержит порты для взаимодействий с картинками в хранилище
type Storage interface {
	// Save сохраняет переданное изображение
	// по заданному пути в хранилище
	Save(file io.Reader, path string) (string, error)

	// Delete удаляет изображение в хранилище по заданному пути
	Delete(path string) error
}

// FirebaseClient содержит порты для взаимодействия с Firebase
type FirebaseClient interface {
	// Send отправляет Push-уведомление пользователю через Firebase
	Send(ctx context.Context,
		message *messaging.Message) (string, error)
}
