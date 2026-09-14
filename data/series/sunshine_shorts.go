package series

import (
	"fmt"
	"gamingdiary/data"
	"gamingdiary/tools/youtube"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type SunshineShorts struct{}

func (SunshineShorts) GetName() string {
	return "sunshine/shorts"
}

func (s SunshineShorts) BuildEntryFromYoutube(index int, video *youtube.Video, entry *data.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Short %02d | %s", index+1, title[0])

	entry.Thumbnail = entry.YoutubeThumbnail
	return nil
}
