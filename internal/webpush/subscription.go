package webpush

import (
	"crypto/ecdh"
	"encoding/base64"
	"time"
)

// Subscription is a Web Push subscription received from a Push Service via a
// Push Manager.
//
// SEE: https://developer.mozilla.org/en-US/docs/Web/API/PushSubscription.
type Subscription struct {
	Endpoint       string           `json:"endpoint"`
	ExpirationTime *time.Time       `json:"expirationTime,omitempty"`
	Keys           SubscriptionKeys `json:"keys"`
}

// PushTarget returns a [PushTarget] for use when pushing messages from an
// Application Server.
func (s *Subscription) PushTarget() (*PushTarget, error) {
	userAgentPublicKey, err := s.Keys.PublicKey()
	if err != nil {
		return nil, err
	}

	authenticationSecret, err := s.Keys.AuthenticationSecret()
	if err != nil {
		return nil, err
	}

	return &PushTarget{
		Endpoint:             s.Endpoint,
		UserAgentPublicKey:   userAgentPublicKey,
		AuthenticationSecret: authenticationSecret,
	}, nil
}

type SubscriptionKeys struct {
	Auth   string `json:"auth"`
	P256DH string `json:"p256dh"`
}

func (s *SubscriptionKeys) PublicKey() (*ecdh.PublicKey, error) {
	publicKey, err := base64.RawURLEncoding.DecodeString(s.P256DH)
	if err != nil {
		return nil, err
	}

	return ecdh.P256().NewPublicKey(publicKey)
}

func (s *SubscriptionKeys) AuthenticationSecret() ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s.Auth)
}
