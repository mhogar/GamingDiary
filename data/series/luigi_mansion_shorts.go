package series

import (
	"app/data"
	"app/tools/youtube"
	"fmt"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type LuigiMansionShorts struct{}

func (LuigiMansionShorts) GetName() string {
	return "luigi_mansion/shorts"
}

func (s LuigiMansionShorts) BuildEntryFromYoutube(index int, video *youtube.Video, entry *data.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Short %d | %s", index+1, title[0])

	entry.Description = strings.SplitN(video.Snippet.Description, "\n", 2)[0]
	entry.Thumbnail = entry.YoutubeThumbnail

	return nil
}
