package heartgold

import (
	"fmt"
	"gamingdiary/data"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

var META_REGEX = regexp.MustCompile(`meta(.+)\.txt$`)

type Parser struct{}

func (Parser) RawFiles(path string) ([]string, error) {
	return filepath.Glob(filepath.Join(path, "download/*.txt"))
}

func (p Parser) ParseEntry(path string) (data.Entry, error) {
	matches := META_REGEX.FindStringSubmatch(path)
	if len(matches) == 0 {
		return data.Entry{}, errors.New("path not match")
	}
	index := matches[1]

	bytes, err := os.ReadFile(path)
	if err != nil {
		return data.Entry{}, errors.Chain(err, "error reading file")
	}
	lines := strings.Split(string(bytes), "\n")

	title := strings.SplitN(lines[2], " | ", 2)

	return data.Entry{
		Title:       fmt.Sprintf("Entry %s | %s", index, data.Capitalize(title[0])),
		Description: lines[7],
		Thumbnail:   fmt.Sprintf("entries/t%s.png", index),
		Video:       fmt.Sprintf("entries/v%s.mp4", index),
		Youtube:     lines[0],
		Groups:      []string{p.selectGroup(index)},
	}, nil
}

func (Parser) selectGroup(index string) string {
	i, err := strconv.ParseInt(index, 10, 32)
	if err != nil {
		return ""
	}

	switch {
	case i >= 1 && i <= 5:
		return "badge1"
	case i >= 6 && i <= 10:
		return "badge2"
	case i >= 11 && i <= 14:
		return "badge3"
	case i >= 15 && i <= 19:
		return "badge4"
	case i >= 20 && i <= 24:
		return "badge5"
	case i >= 25:
		return "badge6"
	default:
		return ""
	}
}
