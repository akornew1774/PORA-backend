// notisend - пакет для отправки сообщений пользователям через Notisend
package notisend

import (
	"pora/internal/dto/responses"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"

	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"os"
)

// SendOtp отправляет запрос на отправку OTP на номер телефона
func SendOtp(phone string, code string) error {

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
