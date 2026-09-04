package shake_it_videos

import (
	"fmt"
	"local/data"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

var META_REGEX = regexp.MustCompile(`meta([0-9]{2})_(a?.+)\.txt$`)

type Parser struct{}

func (Parser) RawFiles(path string) ([]string, error) {
	return filepath.Glob(filepath.Join(path, "youtube/*.txt"))
}

func (p Parser) ParseEntry(path string) (data.Entry, error) {
	matches := META_REGEX.FindStringSubmatch(path)
	if len(matches) == 0 {
		return data.Entry{}, errors.New("path not match")
	}
	index := matches[1]
	area := matches[2]

	bytes, err := os.ReadFile(path)
	if err != nil {
		return data.Entry{}, errors.Chain(err, "error reading file")
	}
	lines := strings.Split(string(bytes), "\n")

	title := strings.SplitN(lines[2], " - ", 2)

	return data.Entry{
		Title:       fmt.Sprintf("%s | %s", p.buildAreaTitle(area), title[1]),
		Description: lines[2],
		Thumbnail:   fmt.Sprintf("t%s_%s.png", index, area),
		Video:       fmt.Sprintf("v%s_%s.mp4", index, area),
		Youtube:     lines[0],
		Groups:      []string{p.selectGroup(area)},
	}, nil
}

func (Parser) buildAreaTitle(area string) string {
	tokens := strings.SplitN(area, "-", 2)
	if len(tokens) < 2 {
		return data.Capitalize(area)
	}
	return fmt.Sprintf("Area %s", area[1:])
}

func (Parser) selectGroup(area string) string {
	tokens := strings.SplitN(area, "-", 2)
	if len(tokens) < 2 {
		return ""
	}
	return "area" + tokens[0][1:]
}
