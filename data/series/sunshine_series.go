package series

import (
	"fmt"
	"gamingdiary/data"
	"gamingdiary/tools/youtube"
	"strings"
)

type SunshineChapters struct{}

func (SunshineChapters) GetName() string {
	return "sunshine/series"
}

func (SunshineChapters) BuildEntryFromYoutube(index string, video *youtube.Video, entry *data.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	entry.Title = fmt.Sprintf("Chapter %s | %s", index, title[0])

	description := strings.SplitN(video.Snippet.Description, "\n", 2)
	entry.Description = description[0]

	return nil
}
