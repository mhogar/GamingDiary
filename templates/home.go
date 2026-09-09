package templates

import (
	"os"
	"text/template"

	"github.com/binarysoupdev/go-extensions/errors"
)

const AUTO_GENERATED_HEADER = "<!-- AUTO GENERATED DO NOT EDIT -->\n"

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

	file.WriteString(AUTO_GENERATED_HEADER)
	return t.Execute(file, data)
}
