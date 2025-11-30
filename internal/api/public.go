package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/AlexGustafsson/grapevine/internal/webpush"
)

type PublicServer struct {
	api API
	mux *http.ServeMux
}

func NewPublicServer(api API) *PublicServer {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/subscriptions/{topic}/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var subscription webpush.Subscription
		if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		topic := r.PathValue("topic")
		id := r.PathValue("id")

		endpointDigest := sha256.Sum256([]byte(subscription.Endpoint))
		expectedID := base64.RawURLEncoding.EncodeToString(endpointDigest[:])

		if id != expectedID {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := api.Subscribe(r.Context(), topic, id, subscription)
		if err == ErrTopicNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("Failed to subscribe", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	mux.HandleFunc("HEAD /api/v1/subscriptions/{topic}/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		// TODO: Proof of ownership of the auth secret used for registering the
		// subscription
		// subscription, err := api.GetSubsription(r.Context(), id)
		// if err != nil {
		// 	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		// 	return
		// }

		_, err := api.GetSubsription(r.Context(), r.PathValue("topic"), id)
		if err == ErrTopicNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err == ErrSubscriptionNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("Failed to get subscription", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("DELETE /api/v1/subscriptions/{topic}/{id}", func(w http.ResponseWriter, r *http.Request) {
		topic := r.PathValue("topic")
		id := r.PathValue("id")

		// TODO: Proof of ownership of the auth secret used for registering the
		// subscription
		// subscription, err := api.GetSubsription(r.Context(), id)
		// if err != nil {
		// 	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		// 	return
		// }

		err := api.Unsubscribe(r.Context(), topic, id)
		if err == ErrTopicNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err == ErrSubscriptionNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("Failed to unsubscribe", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	return &PublicServer{
		mux: mux,
	}
}

func (s *PublicServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
