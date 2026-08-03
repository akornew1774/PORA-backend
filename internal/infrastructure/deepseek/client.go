// deepseek - пакет для отправки запросов к ИИ через Deepseek
package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"pora/internal/config"
	models "pora/internal/dto/ai"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
)

// Client - клиент для выполнения запросов к DeepSeek API
type Client struct {
	cfg        config.DeepseekConfig
	httpClient *http.Client
}

// NewClient создает новый клиент DeepSeek API
func NewClient(cfg config.DeepseekConfig,
	httpClient *http.Client) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Generate отправляет запрос на генерацию ответа модели
func (c *Client) Generate(ctx context.Context,
	request models.AIRequest) (*models.AIResponse, error) {

	body, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error("Ошибка при сериализации запроса: ", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.cfg.BaseURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		logger.Log.Error("Ошибка при создании запроса: ", err)
		return nil, err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+c.cfg.ApiKey,
	)
	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	httpResponse, err := c.httpClient.Do(req)
	if err != nil {
		logger.Log.Error("Ошибка при отправке запроса через HTTP: ", err)
		return nil, err
	}

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		logger.Log.Error("Ошибка при чтении тела ответа: ", err)
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		logger.Log.Errorf(
			"При запросе к Deepseek получен статус %d: %s",
			httpResponse.StatusCode,
			responseBody,
		)
		return nil, errors.ErrorInternal
	}

	var response models.AIResponse

	err = json.Unmarshal(
		responseBody,
		&response,
	)
	if err != nil {
		logger.Log.Error("Ошибка при десериализации ответа: ", err)
		return nil, err
	}

	return &response, nil
}
