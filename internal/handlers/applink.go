// handlers - пакет, содержащий в себе обработчики Api запросов
package handlers

import (
	"net/http"
	"pora/internal/config"
	"pora/internal/dto/responses"

	"github.com/gin-gonic/gin"
)

// AppLinkHandler обрабатывает запросы для
// создания Deep Link на мобильное устройство
type AppLinkHandler struct {
	android  config.AndroidConfig
	deepLink config.DeepLinkConfig
}

// NewAppLinkHandler создает и возвращает новый объект AppLinkHandler
func NewAppLinkHandler(android config.AndroidConfig,
	deepLink config.DeepLinkConfig) *AppLinkHandler {
	return &AppLinkHandler{
		android:  android,
		deepLink: deepLink,
	}
}

// AssetLinks обрабатывает настроику ссылок на приложение для Android
func (h *AppLinkHandler) AssetLinks(c *gin.Context) {

	response := responses.AssetLinksResponse{
		{
			Relation: []string{
				"delegate_permission/common.handle_all_urls",
			},
			Target: responses.AssetTarget{
				Namespace:              "android_app",
				PackageName:            h.android.PackageName,
				Sha256CertFingerprints: []string{h.android.SHA256},
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// OpenInvite переводит пользователя на страницу
// для скачивания приложения (если оно не скачано)
func (h *AppLinkHandler) OpenInvite(c *gin.Context) {
	c.Redirect(
		http.StatusTemporaryRedirect,
		h.deepLink.DownloadLink,
	)
}
