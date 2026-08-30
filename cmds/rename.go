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
	path := cmd.Flags.String("path", "", "search path")
	pattern := cmd.Flags.String("pattern", "(.+)", "the pattern regex")
	replace := cmd.Flags.String("replace", "$", "the replace string")
	confirm := cmd.Flags.Bool("confirm", false, "confirm rename")
	cmd.ParseFlags(args)

	regex, err := regexp.Compile(*pattern)
	if err != nil {
		return errors.Chain(err, "error compiling pattern regex")
	}

	files, err := os.ReadDir(*path)
	if err != nil {
		return errors.Chain(err, "error reading directory")
	}

	count := 0
	for _, file := range files {
		if file.IsDir() || !regex.MatchString(file.Name()) {
			continue
		}
		count++

		matches := regex.FindStringSubmatch(file.Name())
		if len(matches) < 2 {
			return errors.New("pattern regex does not contain a group")
		}

		renamed := strings.Replace(*replace, "$", matches[1], 1)

		if *confirm {
			err := os.Rename(filepath.Join(*path, file.Name()), filepath.Join(*path, renamed))
			if err != nil {
				return errors.Chain(err, "error renaming file")
			}
		} else {
			style.Info.Printf("%s -> %s\n", file.Name(), style.Bold.Sprint(renamed))
		}
	}

	if *confirm {
		style.BoldCreate.Printf("%d Files Renamed\n", count)
		return nil
	}

	return errors.New("rerun with \"-confirm\" to apply")
}
