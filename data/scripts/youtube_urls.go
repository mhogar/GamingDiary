package scripts

import (
	"app/data"
	"fmt"
)

type YoutubeURLs struct{}

func (YoutubeURLs) GetName() string {
	return "youtube/urls"
}

func (YoutubeURLs) Run(_ string, entry *data.Entry) error {
	entry.YoutubeThumbnail = fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", entry.YoutubeId)
	entry.YoutubeVideo = fmt.Sprintf("https://www.youtube.com/watch?v=%s", entry.YoutubeId)
	return nil
}
