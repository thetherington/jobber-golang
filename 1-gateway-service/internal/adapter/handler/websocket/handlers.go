package websocket

import (
	"context"
	"encoding/json"

	"github.com/thetherington/jobber-common/models/chat"
	"github.com/thetherington/jobber-common/models/event"
	"github.com/thetherington/jobber-common/models/order"
)

func (m *Manager) DispatchMessage(cmd string, payload *chat.MessageDocument) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	event := Event{Type: cmd, Payload: json.RawMessage(data)}

	for c := range m.clients {
		if c.username != "" &&
			(c.username == payload.ReceiverUsername || c.username == payload.SenderUsername) {
			c.egress <- event
		}
	}
}

func (m *Manager) DispatchNotification(cmd string, payload *order.Notification) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	event := Event{Type: cmd, Payload: json.RawMessage(data)}
	for c := range m.clients {
		if c.username != "" &&
			(c.username == payload.ReceiverUsername || c.username == payload.SenderUsername) {
			c.egress <- event
		}
	}
}

func (m *Manager) DispatchOrder(cmd string, payload *order.OrderDocument) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	event := Event{Type: cmd, Payload: json.RawMessage(data)}
	for c := range m.clients {
		if c.username != "" &&
			(c.username == payload.SellerUsername || c.username == payload.BuyerUsername) {
			c.egress <- event
		}
	}
}

func (m *Manager) PushLoggedInUsers(users []string) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	ev := Event{Type: event.Online, Payload: json.RawMessage(data)}

	for c := range m.clients {
		c.egress <- ev
	}

	return nil
}

func UpdateClientUsername(e Event, c *Client) error {
	var user string
	if err := json.Unmarshal(e.Payload, &user); err != nil {
		return err
	}

	c.username = user

	return nil
}

func GetLoggedInUsers(_ Event, c *Client) error {
	users, err := c.manager.cache.GetLoggedInUsersFromCache(context.Background())
	if err != nil {
		return err
	}

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	ev := Event{Type: event.Online, Payload: json.RawMessage(data)}

	for c := range c.manager.clients {
		c.egress <- ev
	}

	return nil
}

func LoggedInUsers(e Event, c *Client) error {
	var user string
	if err := json.Unmarshal(e.Payload, &user); err != nil {
		return err
	}

	users, err := c.manager.cache.SaveLoggedInUserToCache(context.Background(), user)
	if err != nil {
		return err
	}

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	ev := Event{Type: event.Online, Payload: json.RawMessage(data)}

	for c := range c.manager.clients {
		c.egress <- ev
	}

	return nil
}

func RemoveLoggedInUser(e Event, c *Client) error {
	var user string
	if err := json.Unmarshal(e.Payload, &user); err != nil {
		return err
	}

	users, err := c.manager.cache.RemoveLoggedInUserFromCache(context.Background(), user)
	if err != nil {
		return err
	}

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	ev := Event{Type: event.Online, Payload: json.RawMessage(data)}

	for c := range c.manager.clients {
		c.egress <- ev
	}

	return nil
}

func SaveUserCategory(e Event, c *Client) error {
	var payload struct {
		Username string `json:"username"`
		Category string `json:"category"`
	}

	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		return err
	}

	return c.manager.cache.SaveUserSelectedCategory(context.Background(), payload.Username, payload.Category)
}
