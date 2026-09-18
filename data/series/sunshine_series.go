package series

import (
	"app/data"
	"app/tools/youtube"
	"fmt"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type SunshineChapters struct {
	seriesBase
}

func (SunshineChapters) GetName() string {
	return "sunshine/series"
}

func (SunshineChapters) BuildEntryFromYoutube(index int, video *youtube.Video, entry *data.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Chapter %02d | %s", index, title[0])

	description := strings.SplitN(video.Snippet.Description, "\n", 2)
	entry.Description = description[0]

	return nil
}
