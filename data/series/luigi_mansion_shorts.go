package series

import (
	"app/data/build"
	"app/tools/youtube"
	"fmt"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type LuigiMansionShorts struct {
	seriesBase
}

func (LuigiMansionShorts) GetName() string {
	return "luigi_mansion/shorts"
}

func (s LuigiMansionShorts) BuildEntryFromYoutube(index, _ int, video *youtube.Video, entry *build.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Short %d | %s", index+1, title[0])

	entry.Description = strings.SplitN(video.Snippet.Description, "\n", 2)[0]
	entry.Thumbnail = entry.YoutubeThumbnail

	return nil
}
