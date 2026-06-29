// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"pora/internal/ports"
)

// FamilyService - объект, содержащий методы для работы с семьями
type FamilyService struct {
	familyRepo ports.FamilyRepository
	listRepo   ports.ListRepository
}

// NewFamilyRepository создает и возвращает новый объект FamilyService
func NewFamilyService(familyRepo ports.FamilyRepository,
	listRepo ports.ListRepository) ports.FamilyService {
	return &FamilyService{
		familyRepo: familyRepo,
		listRepo:   listRepo,
	}
}
