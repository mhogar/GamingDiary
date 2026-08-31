package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
	replace := cmd.Flags.String("replace", "&1", "the replace string")
	num := cmd.Flags.String("num", "0", "iterator starting value and padding")
	step := cmd.Flags.Int("step", 1, "iterator step size")
	groupSize := cmd.Flags.Int("group", 1, "iterator group size")
	confirm := cmd.Flags.Bool("confirm", false, "confirm rename")
	cmd.ParseFlags(args)

	regex, err := regexp.Compile(*pattern)
	if err != nil {
		return errors.Chain(err, "error compiling pattern regex")
	}

	iter64, err := strconv.ParseInt(*num, 10, 16)
	if err != nil {
		return errors.Chain(err, "invalid iterator")
	}

	iter := int(iter64)
	iterFmt := fmt.Sprintf("%%0%dd", len(*num))
	groupCount := 0

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

		renamed := strings.Replace(*replace, "&0", fmt.Sprintf(iterFmt, iter), -1)
		groupCount++

		if groupCount >= *groupSize {
			iter += *step
			groupCount = 0
		}

		for i, match := range matches[1:] {
			renamed = strings.Replace(renamed, fmt.Sprintf("&%d", i+1), match, -1)
		}

		if *confirm {
			newPath := filepath.Join(*path, renamed)

			_, err := os.Stat(newPath)
			if err == nil {
				return errors.Format("new name \"%s\" already exists", newPath)
			}

			err = os.Rename(filepath.Join(*path, file.Name()), newPath)
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
