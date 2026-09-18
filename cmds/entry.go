package cmds

import (
	"app/data"
	"app/data/series"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type EntryCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewEntryCommand() *EntryCommand {
	return &EntryCommand{
		CommandBase: command.NewCommandBase("entry", "create a new entry"),
	}
}

func (cmd *EntryCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd EntryCommand) Run(args []string) error {
	s := cmd.Flags.String("series", "", "name of the series")
	cmd.Flags.Parse(args)

	if *s == "" {
		return errors.New("\"series\" cannot be empty")
	}

	series, err := series.Select(*s)
	if err != nil {
		return err
	}
	style.BoldInfo.Println(*s)

	index, str, err := cmd.calcIndex(series.GetName())
	if err != nil {
		return err
	}
	path := filepath.Join(data.STATIC_PATH, series.GetName(), fmt.Sprintf("entry%s.json", str))

	entry := data.Entry{
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

func (cmd EntryCommand) calcIndex(name string) (int, string, error) {
	entries, err := filepath.Glob(filepath.Join(data.STATIC_PATH, name, data.ENTRY_PATTERN))
	if err != nil {
		return -1, "", errors.Chain(err, "error reading directory")
	}

	if len(entries) == 0 {
		return 0, "0", nil
	}

	matches := data.ENTRY_REGEX.FindStringSubmatch(entries[len(entries)-1])
	if len(matches) < 2 {
		return -1, "", errors.New("invalid entry filename")
	}
	lastIndex := matches[1]

	index, _ := strconv.ParseInt(lastIndex, 10, 16)
	index++

	format := fmt.Sprintf("%%0%dd", len(lastIndex))
	return int(index), fmt.Sprintf(format, index), nil
}
