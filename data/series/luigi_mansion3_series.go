package series

import (
	"app/data"
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

func (s LuigiMansion3Series) BuildEntryFromYoutube(_, videoIndex int, video *youtube.Video, entry *data.Entry) error {
	title := strings.SplitN(video.Snippet.Title, " | ", 2)
	if len(title) < 2 {
		return errors.Format("invalid title")
	}

	if videoIndex == 0 {
		entry.Title = fmt.Sprintf("Prelude | %s", util.Capitalize(title[0]))
	} else {
		entry.Title = fmt.Sprintf("Floor %s | %s", s.calcFloor(videoIndex), util.Capitalize(title[0]))
	}

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
		return "B1"
	case 2:
		return "05"
	default:
		return "XX"
	}
}
