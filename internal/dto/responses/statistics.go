// responses - пакет, содержащий структуры ответов на Api запросы
package responses

import (
	"time"
)

// GetUserProductsResponse - структура для запроса
// на получение всех созданных пользователем товаров
type GetUserProductsResponse struct {
	Items []ItemInfo `json:"items"`
}

// GetPopularProductsResponse - структура для запроса
// на получение популярных продуктов пользователя
type GetPopularProductsResponse struct {
	PopularProducts []PopularProduct `json:"items"`
}

// PopularProduct - модель популярного продукта
type PopularProduct struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	HowOftenEnds int     `json:"how-often-ends"`
	CurrentDay   float64 `json:"current-day"`
}

// GetBriefResponse - структура для ответ на запрос на
// получение всех часто кончающихся товаров пользователя
type GetBriefResponse struct {
	BriefItems []BriefInfo `json:"brief-items"`
}

// BriefInfo - структура, содержащая краткую информацию о товаре
type BriefInfo struct {
	Title   string  `json:"title"`
	Leadind *string `json:"leading,omitempty"`
}

// LoginTimesResponse - структура для ответа на запрос для
// получения времени всех входов пользователя в приложение
type LoginTimesResponse struct {
	LoginTimes []time.Time `json:"login-times"`
}
