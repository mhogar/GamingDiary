package youtube

import (
	"fmt"
	"gamingdiary/data"
	"gamingdiary/tools/youtube"
	"strings"
)

type SunshineChapters struct{}

func (SunshineChapters) BuildEntry(index string, video *youtube.Video, entry *data.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	entry.Title = fmt.Sprintf("Chapter %s | %s", index, title[0])

	description := strings.SplitN(video.Snippet.Description, "\n", 2)
	entry.Description = description[0]

	entry.Video = fmt.Sprintf("v%s.mp4", index)
	entry.Thumbnail = fmt.Sprintf("t%s.png", index)

	return nil
}
