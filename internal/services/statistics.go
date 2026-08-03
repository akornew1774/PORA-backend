// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
	"pora/internal/entities"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"
	"time"

	"github.com/google/uuid"
)

// StatisticsService - объект, содержащий методы для
// предоставления статистики для пользователей
type StatisticsService struct {
	userRepo      ports.UserRepository
	familyRepo    ports.FamilyRepository
	briefItemRepo ports.BriefItemRepository
	userLoginRepo ports.UserLoginRepository
}

// NewStatisticsService создает и возвращает новый объект StatisticsService
func NewStatisticsService(userRepo ports.UserRepository,
	familyRepo ports.FamilyRepository,
	briefItemRepo ports.BriefItemRepository,
	userLoginRepo ports.UserLoginRepository) ports.StatisticsService {
	return &StatisticsService{
		userRepo:      userRepo,
		familyRepo:    familyRepo,
		briefItemRepo: briefItemRepo,
		userLoginRepo: userLoginRepo,
	}
}

// GetUserProducts получает все созданные пользователем товары
func (s *StatisticsService) GetUserProducts(userID uuid.UUID) (
	responses.GetUserProductsResponse, error) {

	var response responses.GetUserProductsResponse
	var items []responses.ItemInfo

	user, err := s.userRepo.FindWithItems(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя: ", err)
		return response, err
	}

	if user == nil {
		logger.Log.Warn("Указанный пользователь не найден")
		return response, errors.ErrorUserNotFound
	}

	for _, item := range user.AddedItems {

		itemInfo := responses.ItemInfo{
			ID: item.ID,

			Name:    item.Name,
			Section: item.Section,

			Quantity: item.Quantity,
			Unit:     item.Unit,

			Priority: item.Priority,
			Urgent:   item.Urgent,

			Checked:         item.Checked,
			RemindEveryDays: item.RemindEveryDays,
		}

		items = append(items, itemInfo)
	}

	response.Items = items
	return response, nil
}

// GetPopularProducts получает популярны продукты пользователя
func (s *StatisticsService) GetPopularProducts(userID uuid.UUID) (
	responses.GetPopularProductsResponse, error) {

	var (
		response        responses.GetPopularProductsResponse
		popularProducts []responses.PopularProduct
		items           []entities.Item
	)

	user, err := s.userRepo.FindWithLists(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя: ", err)
		return response, err
	}

	if user == nil {
		logger.Log.Warn("Указанный пользователь не найден")
		return response, errors.ErrorUserNotFound
	}

	for _, list := range user.Lists {
		for _, item := range list.Items {
			items = append(items, item)
		}
	}

	user, err = s.userRepo.FindWithFamilies(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя: ", err)
		return response, err
	}

	for _, membership := range user.Memberships {

		family, err := s.familyRepo.FindWithLists(membership.FamilyID)
		if err != nil {
			logger.Log.Error("Ошибка при поиске семьи: ", err)
			return response, err
		}

		if family == nil {
			logger.Log.Warn("Указанная семья не найдена")
			return response, errors.ErrorFamilyNotFound
		}

		for _, list := range family.Lists {
			for _, item := range list.Items {
				items = append(items, item)
			}
		}
	}

	itemsByNames := make(map[string][]entities.Item)

	for _, item := range items {

		if item.TimesBought > 0 {
			itemsByNames[item.Name] = append(itemsByNames[item.Name], item)
		}
	}

	for name, items := range itemsByNames {

		var quantity int
		var howOftenEndsSum float64
		var lastTimeBought time.Time

		for _, item := range items {

			quantity = quantity + item.TimesBought

			createdDaysAgo := int(time.Since(item.CreatedAt) / (24 * time.Hour))
			howOftenEnds := float64(createdDaysAgo) / float64(item.TimesBought)

			howOftenEndsSum = howOftenEndsSum + howOftenEnds

			if item.CheckedAt != nil &&
				item.CheckedAt.After(lastTimeBought) {
				lastTimeBought = *item.CheckedAt
			}
		}

		howOftenEnds := int(howOftenEndsSum / float64(len(items)))

		if howOftenEnds == 0 {
			howOftenEnds = 1
		}

		boughtHoursAgo := float64(time.Since(lastTimeBought) / time.Hour)
		currentDay := (boughtHoursAgo / float64(howOftenEnds)) / 24

		popularProduct := responses.PopularProduct{
			Name:         name,
			Quantity:     quantity,
			HowOftenEnds: howOftenEnds,
			CurrentDay:   currentDay,
		}

		popularProducts = append(popularProducts, popularProduct)
	}

	response.PopularProducts = popularProducts

	return response, err
}

// GetBrief получает все быстро кончающиеся товары пользователя
func (s *StatisticsService) GetBrief(userID uuid.UUID) (
	responses.GetBriefResponse, error) {

	var response responses.GetBriefResponse

	briefItems, err := s.briefItemRepo.FindByUserID(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске краткой информации о товарах: ", err)
		return response, err
	}

	for _, briefItem := range briefItems {

		briefItemInfo := responses.BriefInfo{
			Title:   briefItem.Title,
			Leadind: briefItem.Subtitle,
		}

		response.BriefItems = append(response.BriefItems, briefItemInfo)
	}

	return response, nil
}

// SaveBrief сохраняет быстро кончающиеся товары пользователя
func (s *StatisticsService) SaveBrief(userID uuid.UUID,
	req requests.SaveBriefRequest) error {

	err := s.briefItemRepo.DeleteByUserID(userID)
	if err != nil {
		logger.Log.Error("Ошибка при удалении кратких сведений о товарах: ", err)
		return err
	}

	titles := make(map[string]struct{})

	for _, briefItemInfo := range req.BriefItems {

		if _, exists := titles[briefItemInfo.Title]; exists {
			continue
		}

		briefItem := &entities.BriefItem{
			UserID:   userID,
			Title:    briefItemInfo.Title,
			Subtitle: briefItemInfo.Leadind,
		}

		err := s.briefItemRepo.CreateBriefItem(briefItem)
		if err != nil {
			logger.Log.Error("ошибка при сохранении краткой информации о товаре: ", err)
			continue
		}

		titles[briefItem.Title] = struct{}{}
	}

	return nil
}

// GetLoginTimes получает все время входа пользователя в приложение
func (s *StatisticsService) GetLoginTimes(userID uuid.UUID) (
	responses.LoginTimesResponse, error) {

	var response responses.LoginTimesResponse
	var loginTimes []time.Time

	weekAgo := time.Now().UTC().Add(-7 * 24 * time.Hour)

	err := s.userLoginRepo.DeleteOldLogins(weekAgo)
	if err != nil {
		logger.Log.Error("Ошибка при удалении старых записей о входе: ", err)
		return response, err
	}

	logins, err := s.userLoginRepo.FindRecentByUserID(userID, weekAgo)
	if err != nil {
		logger.Log.Error("Ошибка при получении записей о входе: ", err)
		return response, err
	}

	for _, login := range logins {
		loginTimes = append(loginTimes, login.CreatedAt)
	}

	response.LoginTimes = loginTimes
	return response, err
}
