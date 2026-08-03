// config содержит конфигурации
// с данными для работы приложения
package config

import (
	"os"
	"pora/internal/infrastructure/logger"
	"strconv"
)

// Config - объект, содержащий все
// конфигурации для работы приложения
type Config struct {
	Item     ItemConfig
	FireBase FirebaseConfig
	Push     PushConfig
	DeepLink DeepLinkConfig
	Android  AndroidConfig
	Deepseek DeepseekConfig
}

// ItemConfig - конфигурация с нужными значениями для работы с товарами
type ItemConfig struct {
	DefaultReminderTime  string
	UncheckIntervalHours int
}

// FirebaseConfig - конфигурация с нужными для отправки уведомлений через Firebase значениями
type FirebaseConfig struct {
	ProjectID   string
	ClientEmail string
	PrivateKey  string
}

// PushConfig - конфигурация с нужными для отправки Push-уведомлений общими данными
type PushConfig struct {
	NotificationType   string
	NotificationScreen string
}

// DeepLinkConfig - конфигурация для создания ссылок на приложение
type DeepLinkConfig struct {
	Host         string
	DownloadLink string
}

// AndroidConfig - конфигурация для работы сервера с Android (AssetLinks)
type AndroidConfig struct {
	PackageName string
	SHA256      string
}

// DeepseekConfig - конфигурация для отправки запросов к ИИ через DeepseekApi
type DeepseekConfig struct {
	BaseURL string
	ApiKey  string
}

// Load загружает всю информацию для конфигураций
// и возвращает заполненный объект Config
func Load() *Config {
	uncheckIntervalHours, err := strconv.Atoi(
		os.Getenv("UNCHECK_INTERVAL_HOURS"),
	)
	if err != nil {
		logger.Log.Error("Ошибка при конветрации значений .env файла: ", err)
		uncheckIntervalHours = 24
	}

	return &Config{
		Item: ItemConfig{
			DefaultReminderTime:  os.Getenv("DEFAULT_REMINDER_TIME"),
			UncheckIntervalHours: uncheckIntervalHours,
		},

		FireBase: FirebaseConfig{
			ProjectID:   os.Getenv("FIREBASE_PROJECT_ID"),
			ClientEmail: os.Getenv("FIREBASE_CLIENT_EMAIL"),
			PrivateKey:  os.Getenv("FIREBASE_PRIVATE_KEY"),
		},

		Push: PushConfig{
			NotificationType:   os.Getenv("NOTIFICATION_TYPE"),
			NotificationScreen: os.Getenv("NOTIFICATION_SCREEN"),
		},

		DeepLink: DeepLinkConfig{
			Host:         os.Getenv("BASE_URL"),
			DownloadLink: os.Getenv("DOWNLOAD_LINK"),
		},

		Android: AndroidConfig{
			PackageName: os.Getenv("ANDROID_PACKAGE_NAME"),
			SHA256:      os.Getenv("ANDROID_SHA256"),
		},

		Deepseek: DeepseekConfig{
			BaseURL: os.Getenv("DEEPSEEK_BASE_URL"),
			ApiKey:  os.Getenv("DEEPSEEK_API_KEY"),
		},
	}
}
