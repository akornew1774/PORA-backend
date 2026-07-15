// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"pora/internal/dto/requests"
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

	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
	"gopkg.in/mail.v2"
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

		err = s.sendCodeToNumber(phone, code)

	} else {
		err = s.sendCodeToEmail(email, code)
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

// sendCodeToEmail отправляет запрос на отправку OTP на электронную почту
func (s *OtpService) sendCodeToEmail(email string, code string) error {

	message := mail.NewMessage()

	host := os.Getenv("SMTP_HOST")
	port := 587

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	fromName := os.Getenv("SMTP_FROM_NAME")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")

	message.SetHeader(
		"From",
		fmt.Sprintf("%s <%s>", fromName, fromEmail),
	)

	message.SetHeader("To", email)

	message.SetHeader(
		"Subject",
		"Код подтверждения",
	)

	message.SetBody("text/html", fmt.Sprintf(`

		<p>Ваш код подтверждения:</p>

		<h1 style="font-size:32px">%s</h1>

		<p>Никому не сообщайте этот код</p>
	`, code))

	dialer := mail.NewDialer(
		host,
		port,
		username,
		password,
	)

	err := dialer.DialAndSend(message)
	if err != nil {
		logger.Log.Error("Ошибка при отправке сообщения по email: ", err)
		return err
	}

	return nil
}
