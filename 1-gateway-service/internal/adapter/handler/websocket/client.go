package websocket

import (
	"encoding/json"
	"log"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	username   string

	// egress is used to avoid concurrent writes on the websocket connection
	egress chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) PongHandler(pongMsg string) error {
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) ReadMessages() {
	defer func() {
		// cleanup connection
		c.manager.RemoveClient(c)
	}()

	// set ping deadline
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		slog.With("error", err).Error("WSClient: SetReadDeadline")
	}

	// prevent jumbo frames
	c.connection.SetReadLimit(512)

	// handle pings
	c.connection.SetPongHandler(c.PongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			// something with the connection closed that is not graceful
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				slog.With("error", err).Error("WSClient: error reading message")
			}

			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			slog.With("error", err).Error("WSClient: error marshalling event")
			break
		}

		if err := c.manager.RouteEvent(request, c); err != nil {
			log.Println("error handling message:", err)
		}
	}
}

func (c *Client) WriteMessages() {
	defer func() {
		c.manager.RemoveClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					slog.With("error", err).Error("connection closed")
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				slog.With("error", err).Error("WriteMessages: marshal failed")
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				slog.With("error", err).Error("WriteMessages: failed to write to client")
			}

		case <-ticker.C:
			// send a ping to the client
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				slog.With("error", err).Debug("WriteMessages: failed to ping client")
				return
			}
		}
	}
}
