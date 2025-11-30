package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/AlexGustafsson/grapevine/internal/api"
	"github.com/AlexGustafsson/grapevine/internal/state"
	"github.com/AlexGustafsson/grapevine/internal/web"
	"github.com/caarlos0/env/v10"
)

type Config struct {
	BasePath string `env:"BASE_PATH" envDefault:"./config"`
}

func main() {
	var config Config
	if err := env.ParseWithOptions(&config, env.Options{Prefix: "GRAPEVINE_"}); err != nil {
		slog.Error("Failed to parse config", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Migrating state store")
	if err := state.Migrate(config.BasePath); err != nil {
		slog.Error("Failed to migrate state store", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Loading state store")
	store, err := state.Load(config.BasePath)
	if err != nil {
		slog.Error("Failed to load state store", slog.Any("error", err))
		os.Exit(1)
	}

	webPushAPI := &api.WebPushAPI{
		Store: store,
	}

	publicMux := http.NewServeMux()
	publicMux.Handle("/api/v1/", api.NewPublicServer(webPushAPI))
	publicMux.Handle("/", web.NewServer(store))

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
