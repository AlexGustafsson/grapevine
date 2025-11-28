package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

		endpointDigest := sha256.Sum256([]byte(subscription.Endpoint))
		id := base64.RawURLEncoding.EncodeToString(endpointDigest[:])
		fmt.Println(r.PathValue("id"), id)

		if r.PathValue("id") != id {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := api.Subscribe(r.Context(), r.PathValue("topic"), id, subscription); err != nil {
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
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("DELETE /api/v1/subscriptions/{topic}/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		// TODO: Proof of ownership of the auth secret used for registering the
		// subscription
		// subscription, err := api.GetSubsription(r.Context(), id)
		// if err != nil {
		// 	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		// 	return
		// }

		if err := api.Unsubscribe(r.Context(), r.PathValue("topic"), id); err != nil {
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
