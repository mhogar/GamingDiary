package series

import (
	"app/data"
	"app/tools/youtube"
	"fmt"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type ShakeItSeries struct {
	seriesBase
}

func (ShakeItSeries) GetName() string {
	return "shake_it/series"
}

func (s ShakeItSeries) BuildEntryFromYoutube(_, videoIndex int, video *youtube.Video, entry *data.Entry) error {
	area, level := s.calcLevel(videoIndex)

	title := strings.SplitN(video.Snippet.Title, " - ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Area %d-%d | %s", area, level, title[1])

	entry.Description = video.Snippet.Title
	entry.Groups = []string{fmt.Sprintf("area%d", area)}

	return nil
}

func (ShakeItSeries) calcLevel(i int) (int, int) {
	switch {
	case i == 0:
		return 0, 1
	case fitRange(&i, 1, 25):
		return i/5 + 1, i%5 + 1
	case i == 26:
		return 6, 1
	case fitRange(&i, 27, 35):
		return (i/3)*2 + 1, i%3 + 6
	case i == 36:
		return 5, 9
	case fitRange(&i, 37, 42):
		return (i/3)*2 + 2, i%3 + 6
	case i == 43:
		return 4, 9
	case i == 44:
		return 6, 2
	case i == 45:
		return 7, 1
	default:
		return -1, -1
	}
}
