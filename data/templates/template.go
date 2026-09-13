package templates

import (
	"os"
	"text/template"

	"github.com/binarysoupdev/go-extensions/errors"
)

const AUTO_GENERATED_HEADER = "<!-- AUTO GENERATED DO NOT EDIT -->\n"

func renderTemplate(tmpl *template.Template, out string, data any) error {
	file, err := os.Create(out)
	if err != nil {
		return errors.Chain(err, "error creating output file")
	}
	defer file.Close()

	file.WriteString(AUTO_GENERATED_HEADER)
	return tmpl.Execute(file, data)
}
