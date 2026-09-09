package templates

import (
	"os"
	"text/template"

	"github.com/binarysoupdev/go-extensions/errors"
)

type SeriesPage struct {
	Title        string
	Dates        string
	Background   string
	Entries      []Entry
	Theme        string
	Stylesheets  []string
	ResourcePath string
}

type Entry struct {
	Title            string
	Description      string
	Duration         string
	Thumbnail        string
	DefaultThumbnail string
	Video            string
	Youtube          string
	Classes          []string
}

func RenderSeriesPage(path string, data SeriesPage) error {
	t := template.Must(template.ParseFiles("templates/series.gohtml"))

	file, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating output file")
	}
	defer file.Close()

	file.WriteString(AUTO_GENERATED_HEADER)
	return t.Execute(file, data)
}
