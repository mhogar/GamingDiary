package entry_cmd

import (
	"app/data/build"
	"app/data/series"
	"time"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

func (cmd EntryCommand) createEmptyEntry(series series.Series) error {
	index, path, err := cmd.calcNextEntryFromLastFile(series.GetName())
	if err != nil {
		return errors.Chain(err, "error calculating next entry index")
	}

	entry := build.Entry{
		Date: time.Now().Truncate(time.Second),
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
