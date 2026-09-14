package series

import (
	"fmt"
	"gamingdiary/data"
	"gamingdiary/tools/youtube"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type TTYDBattles struct{}

func (TTYDBattles) GetName() string {
	return "ttyd/battles"
}

func (s TTYDBattles) BuildEntryFromYoutube(index int, video *youtube.Video, entry *data.Entry) error {
	chapter, battle := s.calcChapter(index)

	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Battle %d-%d | %s", chapter, battle, title[0])

	entry.Description = strings.Split(video.Snippet.Description, "\n")[0]
	entry.Groups = []string{fmt.Sprintf("chapter%d", chapter)}

	return nil
}

func (s TTYDBattles) calcChapter(i int) (int, int) {
	if shiftRange(&i, 1) {
		return 0, i + 1
	}
	if shiftRange(&i, 4) {
		return 1, i + 1
	}
	if shiftRange(&i, 2) {
		return 2, i + 1
	}
	if shiftRange(&i, 3) {
		return 3, i + 1
	}
	if shiftRange(&i, 3) {
		return 4, i + 1
	}
	if shiftRange(&i, 2) {
		return 5, i
	}
	if shiftRange(&i, 2) {
		return 6, i
	}
	if shiftRange(&i, 1) {
		return 7, i + 1
	}
	if shiftRange(&i, 8) {
		return 8, i + 1
	}
	return 0, 0
}
