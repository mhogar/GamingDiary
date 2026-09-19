package series

import (
	"app/data/build"
	"app/tools/youtube"
	"fmt"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type TTYDShorts struct {
	seriesBase
}

func (TTYDShorts) GetName() string {
	return "ttyd/shorts"
}

func (s TTYDShorts) BuildEntryFromYoutube(index, _ int, video *youtube.Video, entry *build.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Short %03d | %s", index+1, title[0])

	entry.Thumbnail = entry.YoutubeThumbnail
	return nil
}

func (s TTYDShorts) BuildEntryFromVideo(index int, video string, entry *build.Entry) error {
	entry.Title = fmt.Sprintf("Short %03d | ", index+1)
	entry.Thumbnail = ""

	return nil
}
