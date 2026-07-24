// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
	"pora/internal/entities"
	"pora/internal/errors"
	emailInfrastructure "pora/internal/infrastructure/email"
	"pora/internal/infrastructure/logger"
	"pora/internal/infrastructure/notisend"
	"pora/internal/ports"

	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
)

// OTPService - объект, содержащий методы для отправки и подтверждения Otp
type OtpService struct {
	otpRepo      ports.OtpRepository
	userRepo     ports.UserRepository
	deviceRepo   ports.DeviceRepository
	tokenService ports.TokenService
}

// NewOTPService создает и возвращает новый объект OTPService
func NewOTPService(
	otpRepo ports.OtpRepository,
	userRepo ports.UserRepository,
	deviceRepo ports.DeviceRepository,
	tokenService ports.TokenService) ports.OtpService {
	return &OtpService{
		otpRepo:      otpRepo,
		userRepo:     userRepo,
		deviceRepo:   deviceRepo,
		tokenService: tokenService,
	}
}

// SendOtp отправляет Otp-код по номеру телефона и сохраняет его
func (s *OtpService) SendOtp(rawPhone string, email string) error {

	var phone string
	var err error

	if rawPhone == "" && email == "" {
		logger.Log.Warn("В запросе ни указан ни телефон, ни email")
		return errors.ErrorInvalidInput
	}

	if rawPhone != "" && email != "" {
		logger.Log.Warn("В запросе указан и телефон, и email")
		return errors.ErrorInvalidInput
	}

	code, err := s.generateOtp()
	if err != nil {
		logger.Log.Error("Ошибка при создании Otp-кода: ", err)
		return err
	}

	if rawPhone != "" {

		num, err := phonenumbers.Parse(rawPhone, "RU")
		if err != nil || !phonenumbers.IsValidNumber(num) {
			logger.Log.Warn("Некорректный номер телефона: ", rawPhone)
			return errors.ErrorInvalidPhone
		}

		phone = phonenumbers.Format(num, phonenumbers.E164)

		err = notisend.SendOtp(phone, code)

	} else {
		err = emailInfrastructure.SendOtp(email, code)
	}

	if err != nil {
		logger.Log.Error("Ошибка при отправке OTP-кода")
		return errors.ErrorInternal
	}

	newOtp := &entities.Otp{
		Code:  code,
		Phone: phone,
		Email: email,
	}

	err = s.otpRepo.CreateOtp(newOtp)
	if err != nil {
		logger.Log.Error("Ошибка при сохранении OTP-кода")
		return err
	}

	return nil
}

// VerifyOtp сравнивает полученный Otp-код c сохраненным в БД
func (s *OtpService) VerifyOtp(req requests.VerifyOtpRequest) (
	responses.VerifyOtpResponse, error) {

	var (
		otpEntity *entities.Otp
		err       error
		phone     string
	)

	rawPhone := req.Phone
	email := req.Email
	otp := req.OTP

	if rawPhone == "" && email == "" {
		logger.Log.Warn("В запросе ни указан ни телефон, ни email")
		return responses.VerifyOtpResponse{}, errors.ErrorInvalidInput
	}

	if rawPhone != "" && email != "" {
		logger.Log.Warn("В запросе указан и телефон, и email")
		return responses.VerifyOtpResponse{}, errors.ErrorInvalidInput
	}

	if rawPhone != "" {

		num, err := phonenumbers.Parse(rawPhone, "RU")
		if err != nil || !phonenumbers.IsValidNumber(num) {
			return responses.VerifyOtpResponse{},
				errors.ErrorInvalidPhone
		}

		phone = phonenumbers.Format(num, phonenumbers.E164)

		otpEntity, err = s.otpRepo.FindByPhone(phone)

	} else {
		otpEntity, err = s.otpRepo.FindByEmail(email)
	}

	if err != nil {
		logger.Log.Error("Ошибка при поиске OTP по номеру: ", err)
		return responses.VerifyOtpResponse{}, err
	}

	if otpEntity == nil {
		logger.Log.Warn("OTP для данного номера/почты не существует")
		return responses.VerifyOtpResponse{},
			errors.ErrorOTPIncorrect
	}

	if otpEntity.Code != otp || otpEntity.IsUsed {
		return responses.VerifyOtpResponse{},
			errors.ErrorOTPIncorrect
	}

	if time.Now().UTC().After(otpEntity.ExpiresAt) {
		return responses.VerifyOtpResponse{},
			errors.ErrorOTPExpired
	}

	otpEntity.IsUsed = true

	err = s.otpRepo.UpdateOtp(otpEntity)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении Otp: ", err)
		return responses.VerifyOtpResponse{}, err
	}

	var user *entities.User

	if phone != "" {
		user, err = s.userRepo.FindByPhone(phone)
	} else {
		user, err = s.userRepo.FindByEmail(email)
	}

	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя: ", err)
		return responses.VerifyOtpResponse{}, err
	}

	userID := uuid.New()

	var status string
	if user == nil {
		user = &entities.User{
			ID:    userID,
			Phone: phone,
			Email: email,
		}

		err := s.userRepo.CreateUser(user)
		if err != nil {
			logger.Log.Error("Ошибка при создании пользователя: ", err)
			return responses.VerifyOtpResponse{}, err
		}

		status = string(entities.UserStatusNotRegistered)

	} else {
		status = string(user.Status)
	}

	if req.DeviceToken != "" && req.DeviceType != "" {
		err := s.changeUserDevice(user.ID, req.DeviceToken, req.DeviceType)
		if err != nil {
			logger.Log.Warn("Ошибка при обновлении устройства пользователя")
			return responses.VerifyOtpResponse{}, err
		}
	}

	accessToken, refreshToken, err := s.tokenService.CreateTokens(user)
	if err != nil {
		logger.Log.Error("Ошибка при генерации токенов: ", err)
		return responses.VerifyOtpResponse{}, err
	}

	response := responses.VerifyOtpResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Status:       status,
	}
	return response, nil
}

// changeUserDevice создает запись в БД о новом устройстве пользователя
func (s *OtpService) changeUserDevice(userID uuid.UUID,
	deviceToken string, deviceType string) error {

	if deviceType != string(entities.AndroidDevice) &&
		deviceType != string(entities.IOSDevice) {

		logger.Log.Warn("Некорректный тип устройства")
		return errors.ErrorInvalidInput
	}

	device, err := s.deviceRepo.FindByUserID(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске устройства: ", err)
		return err
	}

	if device != nil {

		if device.DeviceToken == deviceToken {
			return nil
		}

		err = s.deviceRepo.DeleteDevice(device)
		if err != nil {
			logger.Log.Error("Ошибка при удалении устройства: ", err)
			return err
		}
	}

	newDevice := &entities.Device{
		UserID:      userID,
		DeviceToken: deviceToken,
		DeviceType:  entities.DeviceType(deviceType),
	}

	err = s.deviceRepo.CreateDevice(newDevice)
	if err != nil {
		logger.Log.Error("Ошибка при создании нового устройства в БД: ", err)
		return err
	}

	return nil
}

// generateOtp генерирует случайный OTP-код
func (s *OtpService) generateOtp() (string, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", errors.ErrorInternal
	}

	return fmt.Sprintf("%06d", number.Int64()), nil
}
