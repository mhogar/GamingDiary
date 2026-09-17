package templates

import (
	"app/data"
	"path/filepath"
	"text/template"
)

var SERIES_TEMPLATE = template.Must(template.ParseFiles(filepath.Join(data.TEMPLATE_PATH, "series.gohtml")))

type SeriesPage struct {
	Title           string
	TotalDuration   string
	Dates           string
	YoutubePlaylist string
	Background      string
	Entries         []Entry
	Theme           string
	Stylesheets     []string
}

type Entry struct {
	Title            string
	Description      string
	Duration         string
	Date             string
	Thumbnail        string
	Video            string
	YoutubeThumbnail string
	YoutubeVideo     string
	Group            string
	Local            bool
}

func RenderSeriesPage(path string, data SeriesPage) error {
	return renderTemplate(SERIES_TEMPLATE, path, data)
}
