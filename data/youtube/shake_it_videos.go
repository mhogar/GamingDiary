package youtube

import (
	"gamingdiary/data"
	"gamingdiary/tools/youtube"
)

type ShakeItVideos struct{}

func (ShakeItVideos) BuildEntry(index string, video *youtube.Video, entry *data.Entry) error {
	entry.Description = video.Snippet.Title
	return nil
}
