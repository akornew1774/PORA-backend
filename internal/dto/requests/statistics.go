// requests - пакет, содержащий структуры запросов по Api
package requests

// SaveBriefRequest - структура для запроса на сохранение
// информации о часто кончающихся товарах
type SaveBriefRequest struct {
	BriefItems []BriefInfo `json:"brief-items" binding:"required"`
}

// BriefInfo - структура, содержащая краткую информацию о товаре
type BriefInfo struct {
	Title   string  `json:"title" binding:"required"`
	Leadind *string `json:"leading"`
}
