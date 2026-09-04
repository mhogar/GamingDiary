package sunshine_chapters

import (
	"fmt"
	"gamingdiary/data"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type Parser struct{}

func (Parser) RawFiles(path string) ([]string, error) {
	return filepath.Glob(filepath.Join(path, "youtube/*.txt"))
}

func (Parser) ParseEntry(path string) (data.Entry, error) {
	index := regexp.MustCompile(`meta(.+)\.txt$`).FindStringSubmatch(path)[1]

	bytes, err := os.ReadFile(path)
	if err != nil {
		return data.Entry{}, errors.Chain(err, "error reading file")
	}
	lines := strings.Split(string(bytes), "\n")

	title := regexp.MustCompile(`^(.+)\s+\|`).FindStringSubmatch(lines[2])[1]

	return data.Entry{
		Title:       fmt.Sprintf("Chapter %s | %s", index, title),
		Description: lines[5],
		Thumbnail:   fmt.Sprintf("t%s.png", index),
		Video:       fmt.Sprintf("v%s.mp4", index),
		Youtube:     lines[0],
	}, nil
}
