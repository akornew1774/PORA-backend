// email - пакет для отправки сообщений пользователям через email
package email

import (
	"fmt"
	"os"

	"pora/internal/infrastructure/logger"

	"github.com/google/uuid"
	"gopkg.in/mail.v2"
)

// SendOtp отправляет запрос на отправку OTP на электронную почту
func SendOtp(email string, code string) error {

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

// SendHelpRequest отправляет запрос на отправку запроса о помощи от админа
func SendHelpRequest(title string,
	msg string, userID uuid.UUID) error {

	message := mail.NewMessage()

	host := os.Getenv("SMTP_HOST")
	port := 587

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	fromName := os.Getenv("SMTP_FROM_NAME")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")

	email := os.Getenv("ADMIN_EMAIL")

	message.SetHeader(
		"From",
		fmt.Sprintf("%s <%s>", fromName, fromEmail),
	)

	message.SetHeader("To", email)

	message.SetHeader(
		"Subject",
		"Вам пришёл запрос на помощь от пользователя",
	)

	message.SetBody("text/html", fmt.Sprintf(`

		<h1 style="font-size:24px">%s</h1>

		<p>%s</p>

		<p>ID пользователя: %s</p>
	`, title, msg, userID.String()))

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
