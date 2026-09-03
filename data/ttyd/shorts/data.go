package ttyd_shorts

import (
	"fmt"
	"local/data"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

var META_REGEX = regexp.MustCompile(`.*s([0-9]+)_c(.+)\.txt$`)

type Parser struct{}

func (Parser) RawFiles(path string) ([]string, error) {
	return filepath.Glob(filepath.Join(path, "meta/*.txt"))
}

func (Parser) ParseEntry(path string) (data.Entry, error) {
	matches := META_REGEX.FindStringSubmatch(path)
	if len(matches) == 0 {
		return data.Entry{}, errors.New("path not match")
	}
	index := matches[1]
	chapter := matches[2]

	bytes, err := os.ReadFile(path)
	if err != nil {
		return data.Entry{}, errors.Chain(err, "error reading file")
	}
	lines := strings.Split(string(bytes), "\n")

	title := strings.SplitN(lines[2], " | ", 2)

	return data.Entry{
		Title:       fmt.Sprintf("Short %s | %s", index, title[0]),
		Description: fmt.Sprintf("From Chapter %s", chapter),
		Thumbnail:   fmt.Sprintf("../chapters/t%s.png", chapter),
		Video:       fmt.Sprintf("s%s_c%s.mp4", index, chapter),
		Youtube:     lines[0],
		Groups:      []string{fmt.Sprintf("chapter%s", chapter[:1])},
	}, nil
}
