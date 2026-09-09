package youtube

import (
	"gamingdiary/data"
	"gamingdiary/tools/youtube"
)

type Meta struct {
	Playlist string `json:"playlist"`
}

type Series interface {
	BuildEntry(index string, video *youtube.Video, entry *data.Entry) error
}
