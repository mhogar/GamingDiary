package entry_cmd

import (
	"app/data"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/binarysoupdev/go-extensions/errors"
)

func (cmd EntryCommand) calcNextEntryFromLastFile(series string) (int, string, error) {
	entries, err := filepath.Glob(filepath.Join(data.STATIC_PATH, series, data.ENTRY_PATTERN))
	if err != nil {
		return -1, "", errors.Chain(err, "error reading directory")
	}

	if len(entries) == 0 {
		index, path := cmd.calcNextEntryFromIndex(series, 0)
		return index, path, nil
	}

	matches := data.ENTRY_REGEX.FindStringSubmatch(entries[len(entries)-1])
	if len(matches) < 2 {
		return -1, "", errors.New("invalid entry filename")
	}
	lastIndex := matches[1]

	index, _ := strconv.ParseInt(lastIndex, 10, 16)
	index++

	format := fmt.Sprintf("%%0%dd", len(lastIndex))
	return int(index), cmd.entryFilepath(series, int(index), format), nil
}

func (cmd EntryCommand) calcNextEntryFromIndex(series string, offset int) (int, string) {
	index := cmd.indexStart + offset
	return index, cmd.entryFilepath(series, index, cmd.indexFormat)
}

func (cmd EntryCommand) entryFilepath(series string, index int, indexFormat string) string {
	return filepath.Join(data.STATIC_PATH, series, fmt.Sprintf("entry%s.json", fmt.Sprintf(indexFormat, index)))
}
