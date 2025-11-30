package api

import (
	"log/slog"
	"net/http"
)

type PrivateServer struct {
	api API
	mux *http.ServeMux
}

func NewPrivateServer(api API) *PrivateServer {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/notifications/{topic}", func(w http.ResponseWriter, r *http.Request) {
		topic := r.PathValue("topic")

		err := api.Push(r.Context(), topic, &Notification{
			Urgency: UrgencyNormal,
			Title:   "Hello, World!",
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
