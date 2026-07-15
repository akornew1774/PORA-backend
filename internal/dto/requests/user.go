// requests - пакет, содержащий структуры запросов по Api
package requests

// UpdateRequest - структура для запроса
// на обновление профиля пользователя
type UpdateUserRequest struct {
	Phone *string `json:"phone,omitempty"`
	Email *string `json:"email,omitempty"`

	Name    *string `json:"name,omitempty"`
	Surname *string `json:"surname,omitempty"`
}

// UpdateDeviceRequest - структура для запроса
// на обновление устройства пользователя
type UpdateDeviceRequest struct {
	DeviceToken string `json:"device-token" binding:"required"`
	DeviceType  string `json:"device-type" binding:"required"`
}
