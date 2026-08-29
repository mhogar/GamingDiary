package cmds

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/got-style/style"
)

type RenameCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewRenameCommand() *RenameCommand {
	return &RenameCommand{
		CommandBase: command.NewCommandBase("rename", "rename files"),
	}
}

func (cmd *RenameCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd RenameCommand) Run(args []string) error {
	s := NewSeriesSelect(cmd.Flags)
	confirm := cmd.Flags.Bool("confirm", false, "confirm rename")
	cmd.ParseFlags(args)

	series, err := s.Select()
	if err != nil {
		return err
	}
	style.BoldInfo.Println(s)

	const PAD = 2
	REGEX := regexp.MustCompile(`.*[^0-9]([0-9]+)(\..+)`)
	path := filepath.Join("data", series)

	files, err := os.ReadDir(path)
	if err != nil {
		return errors.Chain(err, "error reading directory")
	}

	count := 0

	for _, file := range files {
		if file.IsDir() || !REGEX.MatchString(file.Name()) {
			continue
		}

		matches := REGEX.FindStringSubmatch(file.Name())

		index := matches[1]
		ext := matches[2]

		if len(index) >= PAD {
			if !*confirm {
				style.Info.Println(file.Name())
			}
			continue
		}

		renamed := strings.Replace(file.Name(), index+ext, "0"+index+ext, 1)
		if !*confirm {
			style.Create.Printf("%s -> %s\n", file.Name(), renamed)
		}

		if *confirm {
			count++

			err := os.Rename(filepath.Join(path, file.Name()), filepath.Join(path, renamed))
			if err != nil {
				return errors.Chain(err, "error renaming file")
			}
		}
	}

	if *confirm {
		style.BoldCreate.Printf("%d Files Renamed\n", count)
	}
	return nil
}
