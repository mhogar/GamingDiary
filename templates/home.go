package templates

import (
	"os"
	"text/template"

	"github.com/binarysoupdev/go-extensions/errors"
)

type HomePage struct {
	VideoCount    int
	TotalDuration string
	Series        []SeriesHeader
}

type SeriesHeader struct {
	Title         string
	Dates         string
	Description   string
	VideoCount    int
	TotalDuration string
	Thumbnail     string
	Link          string
	Theme         string
}

func RenderHomePage(path string, data HomePage) error {
	t := template.Must(template.ParseFiles("templates/home.gohtml"))

	file, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating output file")
	}
	defer file.Close()

	return t.Execute(file, data)
}
