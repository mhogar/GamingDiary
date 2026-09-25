package series

import (
	"app/data/build"
	"app/tools/youtube"
	"app/util"
	"fmt"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
)

type LuigiMansion3Series struct {
	seriesBase
}

func (LuigiMansion3Series) GetName() string {
	return "luigi_mansion3/series"
}

func (s LuigiMansion3Series) BuildEntryFromVideo(index int, video string, entry *build.Entry) error {
	entry.Title = fmt.Sprintf("%s | ", s.calcFloor(index))
	return nil
}

func (s LuigiMansion3Series) BuildEntryFromYoutube(_, videoIndex int, video *youtube.Video, entry *build.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}
	entry.Title = fmt.Sprintf("%s | %s", s.calcFloor(videoIndex+1), util.Capitalize(title[0]))

	description := strings.Split(video.Snippet.Description, "\n")
	if len(description) < 3 {
		return errors.Format("invalid description")
	}
	entry.Description = description[2]

	return nil
}

func (LuigiMansion3Series) calcFloor(i int) string {
	switch i {
	case 1:
		return "Prelude"
	case 2:
		return "Floor B1"
	case 3:
		return "Floor 05"
	default:
		return "Floor XX"
	}
}
