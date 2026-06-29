// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"pora/internal/ports"
)

type ItemService struct {
	itemRepo ports.ItemRepository
}

// NewItemService создает и возвращает новый объект ItemService
func NewItemService(itemRepo ports.ItemRepository) ports.ItemService {
	return &ItemService{itemRepo: itemRepo}
}
