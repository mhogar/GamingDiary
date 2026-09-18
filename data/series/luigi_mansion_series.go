package series

import (
	"app/data"
	"app/tools/youtube"
	"fmt"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type LuigiMansionSeries struct {
	seriesBase
}

func (LuigiMansionSeries) GetName() string {
	return "luigi_mansion/series"
}

func (s LuigiMansionSeries) BuildEntryFromYoutube(index int, video *youtube.Video, entry *data.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Chapter %02d | %s", index+1, title[0])

	entry.Description = strings.SplitN(video.Snippet.Description, "\n", 2)[0]
	entry.Groups = []string{fmt.Sprintf("area%d", s.calcArea(index))}

	return nil
}

func (LuigiMansionSeries) calcArea(i int) int {
	if shiftRange(&i, 1) {
		return 0
	}
	if shiftRange(&i, 1) {
		return 1
	}
	if shiftRange(&i, 3) {
		return 2
	}
	if shiftRange(&i, 3) {
		return 3
	}
	if shiftRange(&i, 4) {
		return 4
	}
	if shiftRange(&i, 1) {
		return 5
	}
	return 0
}
