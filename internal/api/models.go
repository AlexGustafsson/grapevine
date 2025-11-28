package api

import (
	"github.com/AlexGustafsson/grapevine/internal/webpush"
)

type SubscriptionRequest struct {
	Topic        string               `json:"topic"`
	Subscription webpush.Subscription `json:"subscription"`
}
