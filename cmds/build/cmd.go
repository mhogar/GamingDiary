package build_cmd

import (
	"app/data"
	"app/data/build"
	"app/data/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
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
	p := cmd.Flags.String("public", data.PUBLIC_PATH, "the public path")
	cmd.Flags.BoolVar(&cmd.local, "local", false, "build using local thumbnails and videos")
	cmd.ParseFlags(args)

	if *s == "" {
		return errors.New("\"series\" cannot be empty")
	}
	if *p == "" {
		return errors.New("\"public\" cannot be empty")
	}
	rootPath := filepath.Join(data.STATIC_PATH, "root.json")

	root, err := json.UnmarshalFile[build.Root](rootPath)
	if err != nil {
		return errors.Chain(err, "error reading root file")
	}

	series, err := cmd.selectSeries(*s, root)
	if err != nil {
		return err
	}

	appName, public := cmd.selectPublicPath(*p)
	stat, err := os.Stat(public)
	if err != nil || !stat.IsDir() {
		return errors.Format("invalid public path \"%s\"", public)
	}

	f, err := os.Create(filepath.Join(data.LOGS_PATH, fmt.Sprintf("build-%s.txt", time.Now().Format(time.DateTime))))
	if err != nil {
		return errors.Chain(err, "error creating log file")
	}
	defer f.Close()
	cmd.logger = log.New(f, "", log.Ltime)

	for _, s := range series {
		data, err := cmd.buildSeries(public, s)
		if err == nil {
			root.Series[s] = data
		} else {
			cmd.printError(err.Error())
		}
	}

	if err := cmd.buildRoot(public, appName, root); err != nil {
		cmd.printError(err.Error())
	}

	if err := json.MarshalFilePretty(root, rootPath, "    "); err != nil {
		return errors.Chain(err, "error saving root file")
	}
	return nil
}

func (cmd BuildCommand) selectSeries(name string, root build.Root) ([]string, error) {
	switch name {
	case "root":
		return []string{}, nil
	case "all":
		s := make([]string, 0, len(root.Series))
		for name := range root.Series {
			s = append(s, name)
		}
		return s, nil
	default:
		stat, err := os.Stat(filepath.Join(data.STATIC_PATH, name))
		if err != nil || !stat.IsDir() {
			return nil, errors.Format("invalid series \"%s\"", name)
		}
		return []string{name}, nil
	}
}

func (cmd BuildCommand) selectPublicPath(name string) (string, string) {
	cfg, err := json.UnmarshalFile[config.Config](data.CONFIG_PATH)
	if err != nil {
		return "", name
	}

	path, ok := cfg.PublicPaths[name]
	if ok {
		return name, path
	} else {
		return "", name
	}
}
