package series

import (
	"app/data/build"
	"app/tools/youtube"
	"fmt"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type TTYDSeries struct {
	seriesBase
}

func (TTYDSeries) GetName() string {
	return "ttyd/series"
}

func (s TTYDSeries) BuildEntryFromYoutube(_, videoIndex int, video *youtube.Video, entry *build.Entry) error {
	chapter, sub := s.calcChapter(videoIndex)

	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Chapter %d-%02d | %s", chapter, sub, title[0])

	entry.Description = strings.Split(video.Snippet.Description, "\n")[0]

	entry.Group = fmt.Sprintf("chapter%d", chapter)
	return nil
}

func (s TTYDSeries) calcChapter(i int) (int, int) {
	if shiftRange(&i, 3) {
		return 0, i
	}
	if shiftRange(&i, 10) {
		return 1, i + 1
	}
	if shiftRange(&i, 8) {
		return 2, i + 1
	}
	if shiftRange(&i, 11) {
		return 3, i + 1
	}
	if shiftRange(&i, 11) {
		return 4, i + 1
	}
	if shiftRange(&i, 9) {
		return 5, i
	}
	if shiftRange(&i, 13) {
		return 6, i
	}
	if shiftRange(&i, 7) {
		return 7, i + 1
	}
	if shiftRange(&i, 11) {
		return 8, i + 1
	}
	return 0, 0
}
