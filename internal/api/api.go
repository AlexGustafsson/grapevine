package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/AlexGustafsson/grapevine/internal/state"
	"github.com/AlexGustafsson/grapevine/internal/webpush"
)

var (
	ErrTopicNotFound        = errors.New("topic not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
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
	Store *state.Store
}

// Subscribe implements API.
func (w *WebPushAPI) Subscribe(ctx context.Context, topic string, id string, subscription webpush.Subscription) error {
	err := w.Store.AddSubscription(topic, id, subscription)
	if err == state.ErrTopicNotFound {
		return ErrTopicNotFound
	} else if err != nil {
		return err
	}

	// Could be debounced queue
	go func() {
		if err := w.Store.Save(w.Store.BasePath()); err != nil {
			slog.Error("Failed to save store", slog.Any("error", err))
		}
	}()

	return nil
}

// GetSubsription implements API.
func (w *WebPushAPI) GetSubsription(ctx context.Context, topic string, id string) (webpush.Subscription, error) {
	subscription, err := w.Store.GetSubscription(topic, id)
	switch err {
	case state.ErrSubscriptionNotFound:
		return subscription, ErrSubscriptionNotFound
	case state.ErrTopicNotFound:
		return subscription, ErrTopicNotFound
	default:
		return subscription, err
	}
}

// Unsubscribe implements API.
func (w *WebPushAPI) Unsubscribe(ctx context.Context, topic string, id string) error {
	err := w.Store.DeleteSubscription(topic, id)
	if err == state.ErrTopicNotFound {
		return ErrTopicNotFound
	} else if err == state.ErrSubscriptionNotFound {
		return ErrSubscriptionNotFound
	} else if err != nil {
		return err
	}

	// Could be debounced queue
	go func() {
		if err := w.Store.Save(w.Store.BasePath()); err != nil {
			slog.Error("Failed to save store", slog.Any("error", err))
		}
	}()

	return nil
}

// Push implements API.
func (w *WebPushAPI) Push(ctx context.Context, topic string, notification *Notification) error {
	client, ok := w.Store.Client(topic)
	if !ok {
		return ErrTopicNotFound
	}

	subscriptions, err := w.Store.GetSubscriptions(topic)
	if err == state.ErrTopicNotFound {
		return ErrTopicNotFound
	} else if err != nil {
		return err
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

		t := true
		message := webpush.DeclerativePushMessage{
			WebPush: 8030,
			Notification: webpush.DeclerativePushNotification{
				Title:              notification.Title,
				Navigate:           "https://example.com", // TODO: Get from subscription - must match
				Body:               notification.Body,
				RequireInteraction: &t,
				// TODO: Unknown if actions work
				// Actions: []webpush.DeclerativePushNotificationAction{
				// 	{
				// 		Action:   "Test",
				// 		Title:    "Test T",
				// 		Navigate: "https://example.com/action",
				// 	},
				// },
			},
			// TODO: AppBadge works
			// AppBadge: 1,
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

		err = client.WebPushClient().Push(ctx, target, content, options)
		if err != nil {
			pushErrors = append(pushErrors, err)
		}
	}

	return errors.Join(pushErrors...)
}
