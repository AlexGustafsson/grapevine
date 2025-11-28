package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AlexGustafsson/grapevine/internal/webpush"
)

type Urgency string

const (
	UrgencyVeryLow Urgency = "very-low"
	UrgencyLow     Urgency = "low"
	UrgencyNormal  Urgency = "normal"
	UrgencyHigh    Urgency = "high"
)

type Notification struct {
	TTL     int
	Urgency Urgency
	Title   string
	Body    string
}

type API interface {
	Subscribe(context.Context, string, string, webpush.Subscription) error
	GetSubsription(context.Context, string, string) (webpush.Subscription, error)
	Unsubscribe(context.Context, string, string) error

	Push(context.Context, string, *Notification) error
}

var _ API = (*WebPushAPI)(nil)

type WebPushAPI struct {
	Clients       map[string]webpush.Client
	Subscriptions map[string]map[string]webpush.Subscription
}

// Subscribe implements API.
func (w *WebPushAPI) Subscribe(ctx context.Context, topic string, id string, subscription webpush.Subscription) error {
	subscriptions, ok := w.Subscriptions[topic]
	if !ok {
		subscriptions = make(map[string]webpush.Subscription)
		w.Subscriptions[topic] = subscriptions
	}
	subscriptions[id] = subscription
	return nil
}

// GetSubsription implements API.
func (w *WebPushAPI) GetSubsription(ctx context.Context, topic string, id string) (webpush.Subscription, error) {
	subscriptions, ok := w.Subscriptions[topic]
	if !ok {
		subscriptions = make(map[string]webpush.Subscription)
		w.Subscriptions[topic] = subscriptions
	}

	subscription, ok := subscriptions[id]
	if !ok {
		return webpush.Subscription{}, fmt.Errorf("no such subscription")
	}

	return subscription, nil
}

// Unsubscribe implements API.
func (w *WebPushAPI) Unsubscribe(ctx context.Context, topic string, id string) error {
	subscriptions, ok := w.Subscriptions[topic]
	if !ok {
		return fmt.Errorf("no such subscription")
	}

	_, ok = subscriptions[id]
	if !ok {
		return fmt.Errorf("no such subscription")
	}

	delete(subscriptions, id)
	return nil
}

// Push implements API.
func (w *WebPushAPI) Push(ctx context.Context, topic string, notification *Notification) error {
	client, ok := w.Clients[topic]
	if !ok {
		return fmt.Errorf("no client for topic")
	}

	subscriptions, ok := w.Subscriptions[topic]
	if !ok {
		slog.Warn("Got event for unknown topic", slog.String("topic", topic))
		return nil
	}

	if len(subscriptions) == 0 {
		slog.Warn("Got event for topic without subscriptions", slog.String("topic", topic))
		return nil
	}

	pushErrors := make([]error, 0)
	for _, subscription := range subscriptions {
		target, err := subscription.PushTarget()
		if err != nil {
			return err
		}

		message := webpush.DeclerativePushMessage{
			WebPush: 8030,
			Notification: webpush.DeclerativePushNotification{
				Title:    notification.Title,
				Navigate: "https://example.com", // TODO: Get from subscription - must match
			},
		}

		content, err := json.Marshal(&message)
		if err != nil {
			return err
		}

		options := &webpush.PushOptions{
			TTL:         3600, // TODO
			Urgency:     webpush.Urgency(notification.Urgency),
			ContentType: "application/notification+json",
		}

		// TODO: Loop over all subscriptions
		fmt.Printf("%+v\n", subscription)
		err = client.Push(ctx, target, content, options)
		if err != nil {
			pushErrors = append(pushErrors, err)
		}
	}

	return errors.Join(pushErrors...)
}
