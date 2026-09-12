package series

import (
	"gamingdiary/data"
	"gamingdiary/tools/youtube"
)

type Series interface {
	BuildYoutubeEntry(index string, video *youtube.Video, entry *data.Entry) error
	UpgradeEntry(index string, entry *data.Entry) error
}
