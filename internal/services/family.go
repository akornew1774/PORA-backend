// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	goErrors "errors"
	"fmt"
	"os"
	"pora/internal/config"
	"pora/internal/dto/responses"
	"pora/internal/entities"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
)

// FamilyService - объект, содержащий методы для работы с семьями
type FamilyService struct {
	userRepo       ports.UserRepository
	memberRepo     ports.MemberRepository
	familyRepo     ports.FamilyRepository
	listService    ports.ListService
	deepLinkConfig config.DeepLinkConfig
}

// NewFamilyRepository создает и возвращает новый объект FamilyService
func NewFamilyService(userRepo ports.UserRepository,
	memberRepo ports.MemberRepository,
	familyRepo ports.FamilyRepository,
	listService ports.ListService,
	deepLinkConfig config.DeepLinkConfig) ports.FamilyService {
	return &FamilyService{
		userRepo:       userRepo,
		memberRepo:     memberRepo,
		familyRepo:     familyRepo,
		listService:    listService,
		deepLinkConfig: deepLinkConfig,
	}
}

// GetFamilies получает информацию в всех семьях текущего пользователя
func (s *FamilyService) GetFamilies(userID uuid.UUID) (
	responses.GetFamiliesResponse, error) {

	var response responses.GetFamiliesResponse
	var families []entities.Family

	user, err := s.userRepo.FindWithFamilies(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя: ", err)
		return response, err
	}

	for _, membership := range user.Memberships {
		families = append(families, *membership.Family)
	}

	for _, family := range families {

		owner, members, err := s.convertMembers(&family)
		if err != nil {
			logger.Log.Error("Ошибка при конвертации членов семьи: ", err)
			continue
		}

		familyInfo := responses.FamilyInfo{
			ID:        family.ID,
			Name:      family.Name,
			Owner:     owner,
			Members:   members,
			CreatedAt: family.CreatedAt,
		}

		response.Families = append(response.Families, familyInfo)
	}

	return response, nil
}

// convertMembers приводит всех членов семьи к нужному для response формату
func (s *FamilyService) convertMembers(family *entities.Family) (
	responses.FamilyMemberInfo, []responses.FamilyMemberInfo, error) {

	var owner responses.FamilyMemberInfo
	var members []responses.FamilyMemberInfo

	for _, member := range family.Members {

		user := member.User

		if user == nil {
			logger.Log.Warn("Пользователь, соответствующий члену семьи не найден")
			return owner, members, errors.ErrorUserNotFound
		}

		memberInfo := responses.FamilyMemberInfo{
			UserInfo: responses.UserInfo{
				ID:       user.ID,
				Name:     user.Name,
				Surname:  user.Surname,
				ImageURL: os.Getenv("BASE_URL") + user.ImageURL,
			},
			JoinedAt: member.JoinedAt,
			Color:    string(member.Color),
		}

		if member.Role == entities.OwnerRole {
			owner = memberInfo
		} else {
			members = append(members, memberInfo)
		}
	}

	if owner.ID != family.OwnerID {
		logger.Log.Error("Владелец семьи не найден")
		return owner, members, errors.ErrorInternal
	}

	return owner, members, nil
}

// GetLists получает информацию о списках продуктов конкретной семьи
func (s *FamilyService) GetLists(familyID uuid.UUID) (
	responses.GetListsResponse, error) {

	var response responses.GetListsResponse
	var lists []responses.ListInfo

	family, err := s.familyRepo.FindWithLists(familyID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске семьи: ", err)
		return response, err
	}

	if family == nil {
		logger.Log.Warn("Указанная семья не найдена")
		return response, errors.ErrorFamilyNotFound
	}

	for _, list := range family.Lists {

		sections, err := s.listService.
			GetHighestPrioritySections(&list)

		if err != nil {
			logger.Log.Warn("Ошибка при получении секций списка продуктов: ", err)
			continue
		}

		listInfo := responses.ListInfo{
			ID:        list.ID,
			Name:      list.Name,
			Sections:  sections,
			CreatedAt: list.CreatedAt,
			UpdatedAt: list.UpdatedAt,
		}

		lists = append(lists, listInfo)
	}

	response.Lists = lists
	return response, nil
}

// CreateFamily создает новую семью и возвращает её ID
func (s *FamilyService) CreateFamily(userID uuid.UUID,
	name string) (responses.IDResponse, error) {

	var response responses.IDResponse

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя по ID")
		return response, err
	}

	if user == nil {
		logger.Log.Error("Указанный пользователь не найден")
		return response, errors.ErrorUserNotFound
	}

	familyID := uuid.New()

	family := &entities.Family{
		ID:      familyID,
		OwnerID: userID,
		Name:    name,
	}

	var attemptsLeft int = 10
	var pgErr *pgconn.PgError

	// Если inviteCode не уникален, делается attemptsLeft
	// попыток сделать новый уникальный код
	for attemptsLeft > 0 {

		inviteCode, err := family.GenerateInviteCode()
		if err != nil {
			return response, err
		}

		family.InviteCode = inviteCode

		err = s.familyRepo.CreateFamily(family)
		if err == nil {
			break
		}

		if goErrors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				attemptsLeft--
				continue
			}
		}

		logger.Log.Error("Ошибка при создании семьи: ", err)
		return response, err
	}

	if attemptsLeft == 0 {
		logger.Log.Error("Не удалось создать уникальный код семьи")
		return response, errors.ErrorInternal
	}

	freeColor, err := family.GetFreeColor()
	if err != nil {
		logger.Log.Error("Ошибка при поиске свободного цвета")
		return response, err
	}

	newMember := &entities.FamilyMember{
		UserID:   userID,
		FamilyID: familyID,

		Role:  entities.OwnerRole,
		Color: freeColor,
	}

	err = s.memberRepo.CreateMember(newMember)

	if err != nil {
		logger.Log.Error("Ошибка при создании нового члена семьи")
		return response, err
	}

	response.ID = familyID
	return response, nil
}

// AddMember добавляет к существующей семье еще одного участника
func (s *FamilyService) AddMember(userID uuid.UUID,
	familyCode string) error {

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя: ", err)
		return err
	}

	if user == nil {
		logger.Log.Warn("Указанный пользователь не найден")
		return errors.ErrorUserNotFound
	}

	family, err := s.familyRepo.FindByCode(familyCode)
	if err != nil {
		logger.Log.Error("Ошибка при поиске семьи: ", err)
		return err
	}

	if family == nil {
		logger.Log.Warn("Указанная семья не найдена")
		return errors.ErrorFamilyNotFound
	}

	member, err := s.memberRepo.FindByUserAndFamily(userID, family.ID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске члена семьи: ", err)
		return err
	}

	if member != nil {
		logger.Log.Warn("Пользователь уже является членом семьи")
		return errors.ErrorMemberExists
	}

	freeColor, err := family.GetFreeColor()
	if err != nil {
		logger.Log.Error("Ошибка при поиске свободного члена")
		return err
	}

	newMember := &entities.FamilyMember{
		UserID:   userID,
		FamilyID: family.ID,

		Role:  entities.MemberRole,
		Color: freeColor,
	}

	err = s.memberRepo.CreateMember(newMember)

	if err != nil {
		logger.Log.Error("Ошибка при создании нового члена семьи")
		return err
	}

	return nil
}

// GetFamilyLink получает Link-код семьи
// вместе с ссылкой для вступления в неё
func (s *FamilyService) GetFamilyLink(familyID uuid.UUID) (
	responses.GetFamilyLinkResponse, error) {

	var response responses.GetFamilyLinkResponse

	family, err := s.familyRepo.FindByID(familyID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске семьи: ", err)
		return response, err
	}

	if family == nil {
		logger.Log.Warn("Указанная семья не найдена")
		return response, errors.ErrorFamilyNotFound
	}

	response.LinkCode = family.InviteCode

	linkURL := fmt.Sprintf(
		"%s/api/families/join/%s",
		s.deepLinkConfig.Host,
		family.InviteCode,
	)

	response.LinkURL = linkURL

	return response, nil
}

// DeleteFamily удаляет одну конкретную семью
func (s *FamilyService) DeleteFamily(familyID uuid.UUID) error {

	family, err := s.familyRepo.FindByID(familyID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске семьи: ", err)
		return err
	}

	if family == nil {
		logger.Log.Warn("Указанная семья не найдена")
		return errors.ErrorFamilyNotFound
	}

	err = s.familyRepo.DeleteFamily(family)
	if err != nil {
		logger.Log.Error("Ошибка при удалении семьи: ", err)
	}

	return nil
}
