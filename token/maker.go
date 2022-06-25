package token

import "time"

type Maker interface {
	CreateToken(email string, platform Platform, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
