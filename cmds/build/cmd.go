package build_cmd

import (
	"app/data"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type BuildCommand struct {
	command.CommandBase
	command.FlagCommand

	local  bool
	logger *log.Logger
}

func NewBuildCommand() *BuildCommand {
	return &BuildCommand{
		CommandBase: command.NewCommandBase("build", "Build the app from the templates"),
	}
}

func (cmd *BuildCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd BuildCommand) Run(args []string) error {
	s := cmd.Flags.String("series", "", "name of the series")
	out := cmd.Flags.String("out", data.PUBLIC_PATH, "the destination path")
	cmd.Flags.BoolVar(&cmd.local, "local", false, "build using local thumbnails and videos")
	cmd.ParseFlags(args)

	if *s == "" {
		return errors.New("\"series\" cannot be empty")
	}
	rootPath := filepath.Join(data.STATIC_PATH, "root.json")

	root, err := json.UnmarshalFile[data.Root](rootPath)
	if err != nil {
		return errors.Chain(err, "error reading root file")
	}

	series, err := cmd.selectSeries(*s, root)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(data.LOGS_PATH, fmt.Sprintf("build-%s.txt", time.Now().Format(time.DateTime))))
	if err != nil {
		return errors.Chain(err, "error creating log file")
	}
	defer f.Close()
	cmd.logger = log.New(f, "", log.Ltime)

	for _, s := range series {
		data, err := cmd.buildSeries(*out, s)
		if err == nil {
			root.Series[s] = data
		} else {
			style.Error.Printf("[x] %s\n", err)
		}
	}

	if err := cmd.buildRoot(*out, root); err != nil {
		return err
	}

	if err := json.MarshalFilePretty(root, rootPath, "    "); err != nil {
		return errors.Chain(err, "error saving root file")
	}
	return nil
}
