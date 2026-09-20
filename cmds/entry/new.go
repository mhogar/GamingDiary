package entry_cmd

import (
	"app/data"
	"app/data/build"
	"app/data/series"
	"fmt"
	"time"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

func (cmd EntryCommand) createNewEntry(series series.Series) error {
	index, path, err := cmd.calcNextEntryFromLastFile(series.GetName())
	if err != nil {
		return errors.Chain(err, "error calculating next entry index")
	}
	matches := data.ENTRY_REGEX.FindStringSubmatch(path)

	now := time.Now()
	entry := build.Entry{
		Date:      time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		Thumbnail: fmt.Sprintf("t%s.png", matches[1]),
		Video:     fmt.Sprintf("v%s.mp4", matches[1]),
	}

	if err := series.BuildNewEntry(index, &entry); err != nil {
		return errors.Chain(err, "error building entry")
	}

	if err := json.MarshalFilePretty(entry, path, "    "); err != nil {
		return errors.Chain(err, "error saving entry")
	}

	style.Create.Printf("+ %s\n", path)
	return nil
}
