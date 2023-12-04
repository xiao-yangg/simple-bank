package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Different types of error returned by VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains payload data of token
type Payload struct {
	Username string `json:"username"`
	Role string `json:"role"`
	jwt.RegisteredClaims `json:"payload"` // simplify if using PASETO only
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{}
	payload.ID = tokenID.String()
	payload.Username = username
	payload.Role = role
	payload.IssuedAt = jwt.NewNumericDate(time.Now())
	payload.ExpiresAt = jwt.NewNumericDate(time.Now().Add(duration))
	return payload, nil
}

func (paylad *Payload) Valid() error {
	if time.Now().After(paylad.ExpiresAt.Time) {
		return ErrExpiredToken
	}
	return nil
}