package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type PrivateServer struct {
	api API
	mux *http.ServeMux
}

type NotificationRequest struct {
	TTL     int     `json:"ttl"`
	Urgency Urgency `json:"urgency"`
	Title   string  `json:"title"`
	Body    string  `json:"body"`
}

func NewPrivateServer(api API) *PrivateServer {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/notifications/{topic}", func(w http.ResponseWriter, r *http.Request) {
		topic := r.PathValue("topic")

		var request NotificationRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := api.Push(r.Context(), topic, &Notification{
			TTL:     request.TTL,
			Urgency: request.Urgency,
			Title:   request.Title,
			Body:    request.Body,
		})
		if err == ErrTopicNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("Failed to publish notification", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	return &PrivateServer{
		api: api,
		mux: mux,
	}
}

func (s *PrivateServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
