package entry_cmd

import (
	"app/data"
	"app/data/build"
	"app/data/series"
	"app/tools/ffmpeg"
	"fmt"
	"path/filepath"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

func (cmd EntryCommand) createEntryFromVideo(video string, series series.Series) error {
	index, path, err := cmd.calcNextEntryFromLastFile(series.GetName())
	if err != nil {
		return errors.Chain(err, "error calculating next entry index")
	}

	duration, err := ffmpeg.CalcVideoDuration(video)
	if err != nil {
		return errors.Chain(err, "error calculating video duration")
	}

	entry := build.Entry{
		Date:     today(),
		Duration: duration,
		Video:    filepath.Base(video),
	}

	matches := data.VIDEO_INDEX_REGEX.FindStringSubmatch(video)
	if len(matches) >= 2 {
		entry.Thumbnail = fmt.Sprintf("t%s.png", matches[1])
	}

	if err := series.BuildEntryFromVideo(index, video, &entry); err != nil {
		return errors.Chain(err, "error building entry")
	}

	if err := json.MarshalFilePretty(entry, path, "    "); err != nil {
		return errors.Chain(err, "error saving entry file")
	}

	style.Create.Printf("+ %s\n", path)
	return nil
}
