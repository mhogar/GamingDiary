package series

import (
	"app/data"
	"app/tools/youtube"
	"app/util"
	"fmt"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type HeartgoldSeries struct{}

func (HeartgoldSeries) GetName() string {
	return "heartgold/series"
}

func (s HeartgoldSeries) BuildEntryFromYoutube(index int, video *youtube.Video, entry *data.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Entry %02d | %s", index, util.Capitalize(title[0]))

	description := strings.Split(video.Snippet.Description, "\n")
	if len(description) < 3 {
		return errors.Format("invalid description")
	}
	entry.Description = description[2]

	entry.Groups = []string{fmt.Sprintf("badge%d", s.calcBadge(index))}

	return nil
}

func (HeartgoldSeries) calcBadge(i int) int {
	if shiftRange(&i, 1) {
		return 0
	}
	if shiftRange(&i, 5) {
		return 1
	}
	if shiftRange(&i, 5) {
		return 2
	}
	if shiftRange(&i, 4) {
		return 3
	}
	if shiftRange(&i, 5) {
		return 4
	}
	if shiftRange(&i, 5) {
		return 5
	}
	if shiftRange(&i, 2) {
		return 6
	}
	return 0
}
