// websocket - пакет, содержащий реализацию постоянного
// соединения между сервером и мобильными устройствами
package websocket

import (
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"sync"

	"github.com/gorilla/websocket"
)

// Client описывает одно активное Websocket-соединение
type Client struct {
	conn      *websocket.Conn
	send      chan []byte
	closeOnce sync.Once
}

// NewClient создает и возвращает нового клиента
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 20),
	}
}

// Send помещает сообщение в очередь на отправку для WritePump
func (c *Client) Send(message []byte) error {
	select {
	case c.send <- message:
		return nil
	default:
		logger.Log.Error("Очередь на отправку переполнена")
		return errors.ErrorInternal
	}
}

// ReadPump читает полученные по WS-соединению сообщения
func (c *Client) ReadPump() {

	defer c.Close()

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			logger.Log.Info("Ошика при чтении сообщения через WS: ", err)
			return
		}
	}
}

// WritePump непрерывно отправляет полученные сообщения клиенту
func (c *Client) WritePump() {

	defer c.Close()

	for message := range c.send {

		err := c.conn.WriteMessage(
			websocket.TextMessage,
			message,
		)

		if err != nil {
			logger.Log.Info("Ошибка отправки сообщения по WS: ", err)
			return
		}
	}
}

// Close закрывает активное websocket-соединение
func (c *Client) Close() error {

	var err error

	c.closeOnce.Do(func() {
		close(c.send)
		err = c.conn.Close()
	})

	return err
}
