// websocket - пакет для создания постоянного соединения
// между сервером и мобильным устройством
package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader - объект, превращающий HTTP соединение в WebSocket
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
