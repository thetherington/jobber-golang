package websocket

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/thetherington/jobber-common/models/event"
	"github.com/thetherington/jobber-gateway/internal/adapter/config"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

var websocketUpgrader = websocket.Upgrader{
	CheckOrigin:     checkOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients ClientList
	cache   port.CacheRepository
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager(cache port.CacheRepository) *Manager {
	m := &Manager{
		clients:  make(ClientList),
		cache:    cache,
		handlers: make(map[string]EventHandler),
	}

	m.SetupEventHandlers()
	return m
}

func (m *Manager) RouteEvent(event Event, c *Client) error {
	// Check event type is part of the map
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("RouteEvent: invalid event type")
	}
}

func (m *Manager) SetupEventHandlers() {
	m.handlers[event.UpdateUsername] = UpdateClientUsername
	m.handlers[event.GetLoggedInUsers] = GetLoggedInUsers
	m.handlers[event.LoggedInUsers] = LoggedInUsers
	m.handlers[event.RemoveLoggedInUser] = RemoveLoggedInUser
	m.handlers[event.Category] = SaveUserCategory
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("client connected!")

	// upgrade regular http connection into websocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.With("error", err).Error("failed to upgrade to websocket connection")
		return
	}

	client := NewClient(conn, m)
	client.username = r.URL.Query().Get("username")

	m.AddClient(client)

	// Start client process
	go client.ReadMessages()
	go client.WriteMessages()
}

func (m *Manager) AddClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) RemoveClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case config.Config.App.ClientUrl:
		return true
	default:
		return false
	}
}
