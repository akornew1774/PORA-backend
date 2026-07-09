// ports - пакет, содержащий все порты (интерфейсы)
package ports

import (
	"io"
)

// Storage содержит порты для взаимодействий с картинками в хранилище
type Storage interface {
	// Save сохраняет переданное изображение
	// по заданному пути в хранилище
	Save(file io.Reader, path string) (string, error)

	// Delete удаляет изображение в хранилище по заданному пути
	Delete(path string) error
}
