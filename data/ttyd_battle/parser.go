package ttyd_battle

import (
	"fmt"
	"local/cmds/types"
	"os"
	"regexp"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type Parser struct{}

func (Parser) Parse(path string) (types.Entry, error) {
	index := regexp.MustCompile(`meta(.+)\.txt$`).FindStringSubmatch(path)[1]

	bytes, err := os.ReadFile(path)
	if err != nil {
		return types.Entry{}, errors.Chain(err, "error reading file")
	}
	lines := strings.Split(string(bytes), "\n")

	title := regexp.MustCompile(`^(.+)\s+\|`).FindStringSubmatch(lines[2])[1]

	return types.Entry{
		Filepath:    fmt.Sprintf("entries/entry_%s.json", index),
		Index:       index,
		Title:       fmt.Sprintf("Battle %s | %s", index, title),
		Description: lines[5],
		Thumbnail:   fmt.Sprintf("thumbnails/t%s.png", index),
		Video:       fmt.Sprintf("videos/v%s.mp4", index),
		YouTube:     lines[0],
	}, nil
}
