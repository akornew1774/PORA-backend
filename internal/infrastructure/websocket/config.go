// websocket - пакет для создания постоянного соединения
// между сервером и мобильным устройством
package websocket

import (
	"time"
)

// Конфигурация для websocket
const (
	WriteWait  = 10 * time.Second
	PongWait   = 60 * time.Second
	PingPeriod = 50 * time.Second
)
