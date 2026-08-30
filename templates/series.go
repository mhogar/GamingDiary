package templates

import (
	"os"
	"text/template"

	"github.com/binarysoupdev/go-extensions/errors"
)

type SeriesPage struct {
	Title      string
	Dates      string
	Background string
	Entries    []Entry

	Theme       string
	Stylesheets []string
}

type Entry struct {
	Title       string
	Description string
	Duration    string
	Thumbnail   string
	Video       string

	Classes []string
}

func RenderSeriesPage(path string, data SeriesPage) error {
	t := template.Must(template.ParseFiles("templates/series.gohtml"))

	file, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating output file")
	}
	defer file.Close()

	return t.Execute(file, data)
}
