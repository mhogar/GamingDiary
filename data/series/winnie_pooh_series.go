package series

import (
	"app/data"
	"app/tools/youtube"
	"fmt"
	"regexp"

	"github.com/binarysoupdev/go-extensions/errors"
)

var WINNIE_POOH_TITLE_REGEX = regexp.MustCompile(`^.*- (.+)\(`)

type WinniePoohSeries struct{}

func (WinniePoohSeries) GetName() string {
	return "winnie_pooh/series"
}

func (s WinniePoohSeries) BuildEntryFromYoutube(index int, video *youtube.Video, entry *data.Entry) error {
	title := WINNIE_POOH_TITLE_REGEX.FindStringSubmatch(video.Snippet.Title)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("Episode %d | %s", index+1, title[1])

	return nil
}
