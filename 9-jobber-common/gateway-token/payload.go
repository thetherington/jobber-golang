package token

import (
	"errors"
	"time"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains the payload data of the token
type Payload struct {
	ID       string    `json:"id"`
	IssuedAt time.Time `json:"issued_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(id string) (*Payload, error) {
	payload := &Payload{
		ID:       id,
		IssuedAt: time.Now(),
	}

	return payload, nil
}
