// firebase - пакет для обеспечения взаимодействия
// с Firebase для отправки Push-уведомлений
package firebase

import (
	"context"
	"fmt"
	"pora/internal/config"
	"pora/internal/infrastructure/logger"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"

	"google.golang.org/api/option"
)

// Client отвечает за взаимодействие с Firebase
// для отправки Push-уведомления на устройство
type Client struct {
	messagingClient *messaging.Client
}

// NewClient создает и возвращает новый Client
func NewClient(ctx context.Context,
	cfg config.FirebaseConfig) (*Client, error) {

	credentialsJSON := createCredentialsJSON(cfg)

	app, err := firebase.NewApp(
		ctx,
		nil,
		option.WithCredentialsJSON(
			credentialsJSON,
		),
	)
	if err != nil {
		logger.Log.Error("Ошибка при создании приложения firebase: ", err)
		return nil, err
	}

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		logger.Log.Error("Ошибка при создании Messaging Client: ", err)
		return nil, err
	}

	return &Client{
		messagingClient: messagingClient,
	}, nil
}

// Send отправляет Push-уведомление пользователю через Firebase
func (c *Client) Send(ctx context.Context,
	message *messaging.Message) (string, error) {

	id, err := c.messagingClient.Send(
		ctx,
		message,
	)
	if err != nil {
		return "", err
	}

	return id, nil
}

// createCredentialsJSON cоздает JSON с учетными данными для Firebase
func createCredentialsJSON(cfg config.FirebaseConfig) []byte {

	privateKey := strings.ReplaceAll(
		cfg.PrivateKey,
		`\n`,
		"\n",
	)

	credentialsJSON := fmt.Sprintf(
		`{
			"type": "service_account",
			"project_id": "%s",
			"private_key": "%s",
			"client_email": "%s"
		}`,
		cfg.ProjectID,
		privateKey,
		cfg.ClientEmail,
	)

	return []byte(credentialsJSON)
}
