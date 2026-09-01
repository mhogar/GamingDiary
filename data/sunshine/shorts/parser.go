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
	return filepath.Glob(filepath.Join(path, "download/*.txt"))
}

func (Parser) ParseEntry(path string) (data.Entry, error) {
	matches := regexp.MustCompile(`meta([0-9]+)\.txt$`).FindStringSubmatch(path)
	index := matches[1]

	bytes, err := os.ReadFile(path)
	if err != nil {
		return data.Entry{}, errors.Chain(err, "error reading file")
	}
	lines := strings.Split(string(bytes), "\n")

	title := regexp.MustCompile(`^(.+)\s+\|`).FindStringSubmatch(lines[2])[1]

	return data.Entry{
		Title:     fmt.Sprintf("Short %s | %s", index, title),
		Thumbnail: "../book.png",
		Video:     fmt.Sprintf("s%s.mp4", index),
		YouTube:   lines[0],
	}, nil
}
