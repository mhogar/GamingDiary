package cmds

import (
	"bytes"
	"local/cmds/types"
	"local/data/ttyd_battle"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type Parser interface {
	Parse(index, raw string) (types.Entry, error)
}

//=========================================

type BuildCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewBuildCommand() *BuildCommand {
	return &BuildCommand{
		CommandBase: command.NewCommandBase("build", "build the entry data"),
	}
}

func (cmd *BuildCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd BuildCommand) Run(args []string) error {
	public := cmd.Flags.String("public", "", "the public path")
	//s := NewSeriesSelect(cmd.Flags)
	cmd.ParseFlags(args)

	if *public == "" {
		return errors.New("\"public\" cannot be empty")
	}

	// series, err := s.Select()
	// if err != nil {
	// 	return err
	// }
	// style.BoldInfo.Println(series)

	series := "ttyd_battle"
	dataPath := filepath.Join("data", series)

	data, err := json.UnmarshalFile[SeriesData](filepath.Join(dataPath, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	files, err := filepath.Glob(filepath.Join(dataPath, data.RawFiles))
	if err != nil {
		return errors.Chain(err, "error finding raw files")
	}

	for _, file := range files {
		style.Info.Printf("%s -> ", file)

		entry, err := ttyd_battle.Parser{}.Parse(file)
		if err != nil {
			return errors.Chain(err, "error parsing raw file")
		}

		entry.Duration, err = cmd.calcVideoDuration(filepath.Join(*public, series, entry.Video))
		if err != nil {
			return errors.Chain(err, "error calculating video duration")
		}

		err = json.MarshalFilePretty(entry, filepath.Join(dataPath, entry.Filepath), "  ")
		if err != nil {
			return errors.Chain(err, "error saving entry file")
		}

		style.Create.Println(entry.Filepath)
	}

	return nil
}

func (cmd BuildCommand) calcVideoDuration(path string) (float32, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	exe := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", path)
	exe.Stdout = &stdout
	exe.Stderr = &stderr

	if err := exe.Run(); err != nil {
		return 0, errors.Chain(err, stderr.String())
	}

	f, err := strconv.ParseFloat(stdout.String()[:stdout.Len()-1], 32)
	if err != nil {
		return 0, errors.Chain(err, "error parsing duration")
	}

	return float32(f), nil
}
