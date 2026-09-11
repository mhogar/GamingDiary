package templates

type HomePage struct {
	VideoCount    int
	TotalDuration string
	StartDate     string
	Series        []SeriesHeader
}

type SeriesHeader struct {
	Title         string
	Dates         string
	Description   string
	VideoCount    int
	TotalDuration string
	Thumbnail     string
	SubSeries     map[string]string
	Theme         string
}

func RenderHomePage(path string, data HomePage) error {
	return renderTemplate("templates/home.gohtml", path, data)
}
