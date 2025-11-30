package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/AlexGustafsson/grapevine/internal/api"
	"github.com/AlexGustafsson/grapevine/internal/web"
	"github.com/AlexGustafsson/grapevine/internal/webpush"
)

func main() {
	// TODO: Read from file
	signingKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		slog.Error("Failed to generate keys", slog.Any("error", err))
		os.Exit(1)
	}

	keyExchangeKey, err := signingKey.ECDH()
	if err != nil {
		slog.Error("Failed to generate keys", slog.Any("error", err))
		os.Exit(1)
	}

	fmt.Println("client pk", webpush.NewClient("https://example.com", signingKey, keyExchangeKey).PublicKeyString())

	clients := map[string]webpush.Client{
		"default": webpush.NewClient("https://example.com", signingKey, keyExchangeKey),
	}

	webPushAPI := &api.WebPushAPI{
		Clients:       clients,
		Subscriptions: make(map[string]map[string]webpush.Subscription),
	}

	publicMux := http.NewServeMux()
	publicMux.Handle("/api/v1/", api.NewPublicServer(webPushAPI))
	publicMux.Handle("/", web.NewServer(clients))

	publicServer := &http.Server{
		Addr:    ":8080",
		Handler: publicMux,
	}

	internalMux := http.NewServeMux()
	internalMux.Handle("/api/v1/", api.NewPrivateServer(webPushAPI))

	internalServer := &http.Server{
		Addr:    ":8081",
		Handler: internalMux,
	}

	var wg sync.WaitGroup
	failed := false

	// TODO: Error handling, cancellation
	wg.Go(func() {
		err := publicServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to serve public endpoint", slog.Any("error", err))
			// Close the other server
			internalServer.Close()
			failed = true
		}
	})

	wg.Go(func() {
		err := internalServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to serve internal endpoint", slog.Any("error", err))
			// Close the other server
			publicServer.Close()
			failed = true
		}
	})

	wg.Wait()
	if failed {
		os.Exit(1)
	}
}
