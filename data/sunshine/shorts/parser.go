package sunshine_shorts

import (
	"fmt"
	"local/data"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type Parser struct{}

func (Parser) RawFiles(path string) ([]string, error) {
	return filepath.Glob(filepath.Join(path, "raw/*.txt"))
}

func (Parser) ParseEntry(path string) (data.Entry, error) {
	matches := regexp.MustCompile(`meta([0-9]+)_c([0-9]+)\.txt$`).FindStringSubmatch(path)
	index := matches[1]
	chapter := matches[2]

	bytes, err := os.ReadFile(path)
	if err != nil {
		return data.Entry{}, errors.Chain(err, "error reading file")
	}
	lines := strings.Split(string(bytes), "\n")

	title := regexp.MustCompile(`^(.+)\s+\|`).FindStringSubmatch(lines[2])[1]

	return data.Entry{
		Title:       fmt.Sprintf("Short %s | %s", index, title),
		Description: fmt.Sprintf("From Chapter %s", chapter),
		Thumbnail:   fmt.Sprintf("../chapters/t%s.png", chapter),
		Video:       fmt.Sprintf("s%s_c%s.mp4", index, chapter),
		YouTube:     lines[0],
	}, nil
}
