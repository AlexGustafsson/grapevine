package web

import (
	htmltemplate "html/template"
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
	ManifestPath         string
	ApplicationServerKey string
}
