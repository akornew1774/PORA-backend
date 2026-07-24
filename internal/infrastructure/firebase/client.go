// firebase - пакет для обеспечения взаимодействия
// с Firebase для отправки Push-уведомлений
package firebase

import (
	"context"
	"encoding/json"
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

// credentials описывает JSON с учетными данными Firebase.
type credentials struct {
	Type        string `json:"type"`
	ProjectID   string `json:"project_id"`
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
}

// NewClient создает и возвращает новый Client
func NewClient(ctx context.Context,
	cfg config.FirebaseConfig) (*Client, error) {

	credentialsJSON, err := createCredentialsJSON(cfg)
	if err != nil {
		return nil, err
	}

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

// createCredentialsJSON создает JSON с учетными данными Firebase.
func createCredentialsJSON(
	cfg config.FirebaseConfig,
) ([]byte, error) {

	privateKey := strings.ReplaceAll(
		cfg.PrivateKey,
		`\n`,
		"\n",
	)

	cred := credentials{
		Type:        "service_account",
		ProjectID:   cfg.ProjectID,
		PrivateKey:  privateKey,
		ClientEmail: cfg.ClientEmail,
	}

	return json.Marshal(cred)
}
