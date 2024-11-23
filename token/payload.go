package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = fmt.Errorf("token is expired")
	ErrInvalidToken = fmt.Errorf("token is invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  int64     `json:"issued_at"`
	ExpiredAt int64     `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now().Unix(),
		ExpiredAt: time.Now().Add(duration).Unix(),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(time.Unix(p.ExpiredAt, 0)) {
		return ErrExpiredToken
	}
	return nil
}
