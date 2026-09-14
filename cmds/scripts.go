package cmds

import (
	"app/data"
	"app/data/scripts"
	"fmt"
	"path/filepath"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type ScriptsCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewScriptsCommand() *ScriptsCommand {
	return &ScriptsCommand{
		CommandBase: command.NewCommandBase("scripts", "Run a script on the entry files"),
	}
}

func (cmd *ScriptsCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd ScriptsCommand) Run(args []string) error {
	s := cmd.Flags.String("script", "", "name of the script")
	cmd.ParseFlags(args)

	if *s == "" {
		return errors.New("\"script\" cannot be empty")
	}

	script, err := scripts.Select(*s)
	if err != nil {
		return err
	}
	style.Bold.Println(script.GetName())

	files, err := filepath.Glob(filepath.Join(data.STATIC_DIR, "*", "*", "entry*.json"))
	if err != nil {
		return errors.Chain(err, "error finding entry files")
	}

	for _, file := range files {
		err := cmd.runScript(file, script)
		if err != nil {
			fmt.Println()
			return errors.Chain(err, file)
		}

		fmt.Printf("\r... %s ", style.Create.Sprintf("[+] %s ", file))
	}

	style.Info.Printf("\n%d entries affected\n", len(files))
	return nil
}

func (cmd ScriptsCommand) runScript(path string, script scripts.Scripts) error {
	matches := data.ENTRY_REGEX.FindStringSubmatch(path)
	if len(matches) < 2 {
		return errors.New("invalid entry filename")
	}
	index := matches[1]

	entry, err := json.UnmarshalFile[data.Entry](path)
	if err != nil {
		return errors.Chain(err, "error reading entry file")
	}

	if err := script.Run(index, &entry); err != nil {
		return errors.Chain(err, "error running script")
	}

	if err := json.MarshalFilePretty(entry, path, "    "); err != nil {
		return errors.Chain(err, "error saving entry file")
	}

	return nil
}
