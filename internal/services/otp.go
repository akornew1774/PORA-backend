// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"pora/internal/dto/responses"
	"pora/internal/entities"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"

	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/nyaruka/phonenumbers"
)

// OTPService - объект, содержащий методы для отправки и подтверждения Otp
type OtpService struct {
	otpRepo      ports.OtpRepository
	userRepo     ports.UserRepository
	tokenService ports.TokenService
}

// NewOTPService создает и возвращает новый объект OTPService
func NewOTPService(
	otpRepo ports.OtpRepository,
	userRepo ports.UserRepository,
	tokenService ports.TokenService) ports.OtpService {
	return &OtpService{
		otpRepo:      otpRepo,
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// SendOtp отправляет Otp-код по номеру телефона и сохраняет его
func (s *OtpService) SendOtp(rawPhone string) error {

	num, err := phonenumbers.Parse(rawPhone, "RU")
	if err != nil || !phonenumbers.IsValidNumber(num) {
		logger.Log.Warn("Некорректный номер телефона: ", rawPhone)
		return errors.ErrorInvalidPhone
	}

	phone := phonenumbers.Format(num, phonenumbers.E164)
	code, err := s.generateOtp()
	if err != nil {
		logger.Log.Error("Ошибка при создании Otp-кода: ", err)
		return err
	}

	err = s.sendCodeToNumber(phone, code)
	if err != nil {
		logger.Log.Error("Ошибка при отправке OTP-кода")
		return errors.ErrorInternal
	}

	newOtp := &entities.Otp{
		Code:  code,
		Phone: phone,
	}

	err = s.otpRepo.CreateOtp(newOtp)
	if err != nil {
		logger.Log.Error("Ошибка при сохранении OTP-кода")
		return err
	}

	return nil
}

// VerifyOtp сравнивает полученный Otp-код c сохраненным в БД
func (s *OtpService) VerifyOtp(rawPhone string, otp string) (
	responses.VerifyOtpResponse, error) {

	num, err := phonenumbers.Parse(rawPhone, "RU")
	if err != nil || !phonenumbers.IsValidNumber(num) {
		return responses.VerifyOtpResponse{},
			errors.ErrorInvalidPhone
	}

	phone := phonenumbers.Format(num, phonenumbers.E164)

	otpEntity, err := s.otpRepo.FindByPhone(phone)
	if err != nil {
		return responses.VerifyOtpResponse{}, err
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

	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя по номеру: ", err)
		return responses.VerifyOtpResponse{}, err
	}

	var status string
	if user == nil {
		user = &entities.User{
			Phone: phone,
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

// generateOtp генерирует случайный OTP-код
func (s *OtpService) generateOtp() (string, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", errors.ErrorInternal
	}

	return fmt.Sprintf("%06d", number.Int64()), nil
}

// sendCodeToNumber отправляет запрос на отправку OTP на номер телефона
func (s *OtpService) sendCodeToNumber(phone string, code string) error {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	project := os.Getenv("PROJECT_NAME")
	apiKey := os.Getenv("NOTISEND_API_KEY")
	url := os.Getenv("NOTISEND_URL")

	writer.WriteField("project", project)
	writer.WriteField("recipients", phone)
	writer.WriteField("message", code)
	writer.WriteField("apikey", apiKey)

	writer.Close()

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		body,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logger.Log.Error("Ошибка при отправке запроса на отправку OTP: ", err)
		return err
	}
	defer response.Body.Close()

	var result responses.NotisendResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		logger.Log.Error("Ошибка при декодировании ответа: ", err)
		return err
	}

	if result.Status != "success" {
		logger.Log.Error("Не удалость отправить OTP-код, Status: ", result.Status)
		return errors.ErrorInternal
	}

	return nil
}
