package scripts

import (
	"app/data/build"
	"fmt"
)

type YoutubeURLs struct{}

func (YoutubeURLs) GetName() string {
	return "youtube/urls"
}

func (YoutubeURLs) Run(_ string, entry *build.Entry) error {
	entry.YoutubeThumbnail = fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", entry.YoutubeId)
	entry.YoutubeVideo = fmt.Sprintf("https://www.youtube.com/watch?v=%s", entry.YoutubeId)
	return nil
}
