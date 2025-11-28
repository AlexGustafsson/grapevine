package web

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AlexGustafsson/grapevine/internal/webpush"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer(clients map[string]webpush.Client) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		client, ok := clients["grapevine"]
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		indexhtmlTemplate.Execute(w, IndexData{
			ManifestPath:         "/manifest.json",
			ApplicationServerKey: client.PublicKeyString(),
		})
	})

	mux.HandleFunc("GET /index.html", func(w http.ResponseWriter, r *http.Request) {
		client, ok := clients["grapevine"]
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		indexhtmlTemplate.Execute(w, IndexData{
			ManifestPath:         "/manifest.json",
			ApplicationServerKey: client.PublicKeyString(),
		})
	})

	mux.HandleFunc("GET /manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		manifestjsonTemplate.Execute(w, ManifestData{
			ID:        "grapevine",
			ShortName: "Grapevine",
			Name:      "Grapevine",
			Icon:      "/grapevine.png",
			StartURL:  "/",
		})
	})

	mux.HandleFunc("GET /grapevine.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		f, _ := os.OpenFile("grapevine.png", os.O_RDONLY, 0)
		defer f.Close()
		io.Copy(w, f)
	})

	mux.HandleFunc("GET /{name}/index.html", func(w http.ResponseWriter, r *http.Request) {
		client, ok := clients["grapevine"]
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		indexhtmlTemplate.Execute(w, IndexData{
			ManifestPath:         fmt.Sprintf("/%s/manifest.json", r.PathValue("name")),
			ApplicationServerKey: client.PublicKeyString(),
		})
	})

	mux.HandleFunc("GET /{name}/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		name := r.PathValue("name")
		manifestjsonTemplate.Execute(w, ManifestData{
			ID:        name,
			ShortName: name,
			Name:      name,
			Icon:      fmt.Sprintf("/%s/icon.png", name),
			StartURL:  fmt.Sprintf("/%s", name),
		})
	})

	mux.HandleFunc("GET /{name}/icon.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		name := r.PathValue("name")
		f, _ := os.OpenFile(fmt.Sprintf("%s.png", name), os.O_RDONLY, 0) // Lord forgive me for my sins
		defer f.Close()
		io.Copy(w, f)
	})

	return &Server{
		mux: mux,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
