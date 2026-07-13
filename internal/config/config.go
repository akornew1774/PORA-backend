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
	DeepLink DeepLinkConfig
	Android  AndroidConfig
}

// ItemConfig - конфигурация с нужными значениями для работы с товарами
type ItemConfig struct {
	DefaultReminderTime  string
	UncheckIntervalHours int
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

		DeepLink: DeepLinkConfig{
			Host:         os.Getenv("BASE_URL"),
			DownloadLink: os.Getenv("DOWNLOAD_LINK"),
		},

		Android: AndroidConfig{
			PackageName: os.Getenv("ANDROID_PACKAGE_NAME"),
			SHA256:      os.Getenv("ANDROID_SHA256"),
		},
	}
}
