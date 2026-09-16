package templates

import (
	"app/data"
	"path/filepath"
	"text/template"
)

var ROOT_TEMPLATE = template.Must(template.ParseFiles(filepath.Join(data.TEMPLATE_PATH, "root.gohtml")))

type RootPage struct {
	Background    string
	Logo          string
	VideoCount    int
	TotalDuration string
	Dates         string
	Series        []SeriesHeader
}

type SeriesHeader struct {
	Index          int
	Title          string
	Dates          string
	Description    string
	VideoCount     int
	TotalDuration  string
	Thumbnail      string
	SubSeriesLinks []SubSeriesLink
	Theme          string
}

type SubSeriesLink struct {
	Title     string
	Link      string
	Separator string
}

func RenderRootPage(path string, data RootPage) error {
	return renderTemplate(ROOT_TEMPLATE, path, data)
}
