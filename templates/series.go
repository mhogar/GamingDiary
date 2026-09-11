package templates

type SeriesPage struct {
	Title         string
	TotalDuration string
	Dates         string
	Background    string
	Entries       []Entry
	Theme         string
	Stylesheets   []string
	ResourcePath  string
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
}

func RenderSeriesPage(path string, data SeriesPage) error {
	return renderTemplate("templates/series.gohtml", path, data)
}
