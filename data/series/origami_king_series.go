package series

import (
	"app/data"
	"app/tools/youtube"
	"app/util"
	"fmt"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type OrigamiKingSeries struct{}

func (OrigamiKingSeries) GetName() string {
	return "origami_king/series"
}

func (s OrigamiKingSeries) BuildEntryFromYoutube(index int, video *youtube.Video, entry *data.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Chapter %02d | %s", index, util.Capitalize(title[0]))

	description := strings.Split(video.Snippet.Description, "\n")
	if len(description) < 3 {
		return errors.Format("invalid description")
	}
	entry.Description = description[2]

	entry.Group = fmt.Sprintf("chapter%d", s.calcChapter(index))

	return nil
}

func (OrigamiKingSeries) calcChapter(i int) int {
	if shiftRange(&i, 1) {
		return 0
	}
	if shiftRange(&i, 9) {
		return 1
	}
	if shiftRange(&i, 2) {
		return 0
	}
	if shiftRange(&i, 5) {
		return 2
	}
	return 0
}
