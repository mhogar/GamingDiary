package ttyd_chapters

import (
	"fmt"
	"gamingdiary/data"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

var CHAPTER_REGEX = regexp.MustCompile(`^Chapter ([1-8])`)

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

type Upgrader struct{}

func (Upgrader) UpgradeEntry(entry *data.Entry) error {
	matches := CHAPTER_REGEX.FindStringSubmatch(entry.Title)
	if len(matches) < 2 {
		return nil
	}

	entry.Groups = []string{fmt.Sprintf("chapter%s", matches[1])}
	return nil
}
