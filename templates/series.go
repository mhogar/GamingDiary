package templates

import "text/template"

var SERIES_TEMPLATE = template.Must(template.ParseFiles("templates/series.gohtml"))

type SeriesPage struct {
	Title         string
	TotalDuration string
	Dates         string
	Background    string
	Entries       []Entry
	Theme         string
	Stylesheets   []string
}

type Entry struct {
	Title            string
	Description      string
	Duration         string
	Date             string
	Thumbnail        string
	DefaultThumbnail string
	Video            string
	Youtube          string
	Classes          []string
	Local            bool
}

func RenderSeriesPage(path string, data SeriesPage) error {
	return renderTemplate(SERIES_TEMPLATE, path, data)
}
