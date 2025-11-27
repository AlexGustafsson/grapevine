package main

import (
	"fmt"
	htmltemplate "html/template"
	"io"
	"net/http"
	"os"
	texttemplate "text/template"

	_ "embed" // Embed templates
)

//go:embed manifest.json.gotmpl
var manifestjson string

var manifestjsonTemplate = texttemplate.Must(texttemplate.New("name").Parse(manifestjson))

type ManifestData struct {
	ID        string
	Name      string
	ShortName string
	Icon      string
	StartURL  string
}

//go:embed index.html.gotmpl
var indexhtml string

var indexhtmlTemplate = htmltemplate.Must(htmltemplate.New("name").Parse(indexhtml))

type IndexData struct {
	ManifestPath string
}

// NOTE: This is a PoC, there's a lot of unsafe code in this file
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		indexhtmlTemplate.Execute(w, IndexData{
			ManifestPath: "/manifest.json",
		})
	})

	mux.HandleFunc("GET /index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		indexhtmlTemplate.Execute(w, IndexData{
			ManifestPath: "/manifest.json",
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
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		name := r.PathValue("name")
		indexhtmlTemplate.Execute(w, IndexData{
			ManifestPath: fmt.Sprintf("/%s/manifest.json", name),
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

	http.ListenAndServe(":8080", mux)
}
