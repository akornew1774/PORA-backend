// responses - пакет, содержащий структуры ответов на Api запросы
package responses

// GetUserProductsResponse - структура для запроса
// на получение всех созданных пользователем товаров
type GetUserProductsResponse struct {
	Items []ItemInfo `json:"items"`
}

// GetBriefResponse - структура для запроса на
// получение всех часто кончающихся товаров пользователя
type GetBriefResponse struct {
	BriefItems []BriefInfo `json:"brief-items"`
}

// BriefInfo - структура, содержащая краткую информацию о товаре
type BriefInfo struct {
	Title   string  `json:"title"`
	Leadind *string `json:"leading,omitempty"`
}
