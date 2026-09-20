package entry_cmd

import (
	"app/data/series"
	"fmt"
	"strconv"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/got-style/style"
)

type EntryCommand struct {
	command.CommandBase
	command.FlagCommand

	indexStart  int
	indexFormat string
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
	index := cmd.Flags.String("index", "00", "starting index and padding")
	video := cmd.Flags.String("video", "", "create entry from an exiting video")
	youtube := cmd.Flags.String("youtube", "", "create entries from cached youtube data")
	cmd.Flags.Parse(args)

	if *s == "" {
		return errors.New("\"series\" cannot be empty")
	}

	if err := cmd.parseIndexFlag(*index); err != nil {
		return err
	}

	series, err := series.Select(*s)
	if err != nil {
		return err
	}
	style.BoldInfo.Println(*s)

	switch {
	case *video != "":
		return cmd.createEntryFromVideo(*video, series)
	case *youtube != "":
		return cmd.createEntriesFromYoutube(*youtube, series)
	default:
		return cmd.createNewEntry(series)
	}
}

func (cmd *EntryCommand) parseIndexFlag(index string) error {
	i64, err := strconv.ParseInt(index, 10, 16)
	if err != nil {
		return errors.Chain(err, "invalid index")
	}

	cmd.indexStart = int(i64)
	cmd.indexFormat = fmt.Sprintf("%%0%dd", len(index))
	return nil
}
