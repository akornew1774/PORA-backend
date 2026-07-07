// requests - пакет, содержащий структуры запросов по Api
package requests

// FamilyLinkRequest структура для запроса для
// получения ссылки на указанную семью
type FamilyLinkRequest struct {
	FamilyID string `form:"family-id" binding:"required"`
}
