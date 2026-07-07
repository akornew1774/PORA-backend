// config содержит конфигурации
// с данными для работы приложения
package config

import "os"

// Config - объект, содержащий все
// конфигурации для работы приложения
type Config struct {
	DeepLink DeepLinkConfig
	Android  AndroidConfig
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
	return &Config{
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
