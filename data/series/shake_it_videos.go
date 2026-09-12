package series

import (
	"gamingdiary/data"
	"gamingdiary/tools/youtube"
)

type ShakeItVideos struct{}

func (ShakeItVideos) BuildYoutubeEntry(index string, video *youtube.Video, entry *data.Entry) error {
	entry.Description = video.Snippet.Title
	return nil
}

func (ShakeItVideos) UpgradeEntry(index string, entry *data.Entry) error {
	return nil
}
