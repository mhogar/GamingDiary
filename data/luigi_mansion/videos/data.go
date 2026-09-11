package luigi_mansion_videos

import (
	"fmt"
	"gamingdiary/data"
	"gamingdiary/util"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

var META_REGEX = regexp.MustCompile(`meta([0-9]{2})\.txt$`)

type Parser struct{}

func (Parser) RawFiles(path string) ([]string, error) {
	return filepath.Glob(filepath.Join(path, "meta/*.txt"))
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
		Title:       fmt.Sprintf("Entry %s | %s", index, util.Capitalize(title[0])),
		Description: lines[5],
		Thumbnail:   fmt.Sprintf("t%s.png", index),
		Video:       fmt.Sprintf("v%s.mp4", index),
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
	case i == 2:
		return "area1"
	case i >= 3 && i <= 5:
		return "area2"
	case i >= 6 && i <= 8:
		return "area3"
	case i >= 9 && i <= 12:
		return "area4"
	case i == 13:
		return "area5"
	default:
		return ""
	}
}
