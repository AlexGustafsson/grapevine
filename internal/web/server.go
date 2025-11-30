package web

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	texttemplate "text/template"

	"github.com/AlexGustafsson/grapevine/internal/state"
)

//go:embed public
var public embed.FS

//go:embed public/index.html
var index string

type IndexData struct {
	ManifestPath         string
	ApplicationServerKey string
	Topic                string
}

//go:embed manifest.json.gotmpl
var manifest string

type ManifestData struct {
	ID        string
	Name      string
	ShortName string
	Icon      string
	StartURL  string
}

type Server struct {
	mux *http.ServeMux
}

func NewServer(store *state.Store) *Server {
	indexTemplate, err := texttemplate.New("").Parse(index)
	if err != nil {
		panic(err)
	}

	manifestTemplate, err := texttemplate.New("").Parse(manifest)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		client, ok := store.Client("default")
		if !ok {
			// NOTE: The default client should always exist
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = indexTemplate.Execute(w, IndexData{
			ManifestPath:         "/topics/default/manifest.json",
			ApplicationServerKey: client.WebPushClient().PublicKeyString(),
			Topic:                "default",
		})
		if err != nil {
			slog.Error("Failed to render index.html", slog.Any("error", err))
			return
		}
	})

	mux.HandleFunc("GET /topics/{topic}", func(w http.ResponseWriter, r *http.Request) {
		topic := r.PathValue("topic")

		client, ok := store.Client(topic)
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = indexTemplate.Execute(w, IndexData{
			ManifestPath:         fmt.Sprintf("/topics/%s/manifest.json", url.PathEscape(topic)),
			ApplicationServerKey: client.WebPushClient().PublicKeyString(),
			Topic:                url.PathEscape(topic),
		})
		if err != nil {
			slog.Error("Failed to render index.html", slog.Any("error", err))
			return
		}
	})

	mux.HandleFunc("GET /topics/{topic}/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		topic := r.PathValue("topic")

		client, ok := store.Client(topic)
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = manifestTemplate.Execute(w, ManifestData{
			ID:        client.Topic(),
			ShortName: client.ShortName(),
			Name:      client.Name(),
			Icon:      fmt.Sprintf("/topics/%s/icon.png", topic),
			StartURL:  fmt.Sprintf("/topics/%s", topic),
		})
		if err != nil {
			slog.Error("Failed to render manifest.json", slog.Any("error", err))
			return
		}
	})

	mux.HandleFunc("GET /topics/{topic}/icon.png", func(w http.ResponseWriter, r *http.Request) {
		topic := r.PathValue("topic")

		_, ok := store.Client(topic)
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		// TODO - get from client config
		f, _ := os.OpenFile(fmt.Sprintf("%s.png", topic), os.O_RDONLY, 0) // Lord forgive me for my sins
		defer f.Close()
		io.Copy(w, f)
	})

	// Serve public assets
	assets, err := fs.Sub(public, "public")
	if err != nil {
		panic(err)
	}

	mux.Handle("GET /assets/", http.FileServerFS(assets))

	return &Server{
		mux: mux,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
