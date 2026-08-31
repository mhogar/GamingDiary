package cmds

import (
	"bytes"
	"fmt"
	"local/cmds/types"
	"local/data/sunshine"
	ttyd_battle "local/data/ttyd/battles"
	ttyd_chapter "local/data/ttyd/chapters"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type Parser interface {
	Parse(raw string) (types.Entry, error)
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
	src := cmd.Flags.String("src", "", "the source path")
	s := NewSeriesSelect(cmd.Flags)
	cmd.ParseFlags(args)

	if *src == "" {
		return errors.New("\"src\" cannot be empty")
	}

	series, err := s.Select()
	if err != nil {
		return err
	}
	style.BoldInfo.Println(series)

	dataPath := filepath.Join("data", series)
	indexPath := filepath.Join(dataPath, "index.json")

	data, err := json.UnmarshalFile[SeriesData](indexPath)
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	data.VideoCount = 0
	data.TotalDuration = 0

	files, err := filepath.Glob(filepath.Join(dataPath, data.RawFiles))
	if err != nil {
		return errors.Chain(err, "error finding raw files")
	}

	err = os.MkdirAll(filepath.Join(dataPath, filepath.Dir(data.Entries)), 0755)
	if err != nil {
		return errors.Chain(err, "error creating entry directory")
	}

	for _, file := range files {
		style.Info.Printf("%s -> ", file)
		data.VideoCount++

		entry, err := cmd.selectParser(series).Parse(file)
		if err != nil {
			return errors.Chain(err, "error parsing raw file")
		}

		entry.Duration, err = cmd.calcVideoDuration(filepath.Join(*src, series, entry.Video))
		if err != nil {
			return errors.Chain(err, "error calculating video duration")
		}
		data.TotalDuration += entry.Duration

		out := filepath.Join(dataPath, strings.Replace(data.Entries, "*", entry.Index, 1))

		err = json.MarshalFilePretty(entry, out, "  ")
		if err != nil {
			return errors.Chain(err, "error saving entry file")
		}

		style.Create.Println(out)
	}

	_ = json.MarshalFilePretty(data, indexPath, "    ")
	return nil
}

func (cmd BuildCommand) selectParser(series string) Parser {
	switch series {
	case "sunshine":
		return sunshine.Parser{}
	case "ttyd/chapters":
		return ttyd_chapter.Parser{}
	case "ttyd/battles":
		return ttyd_battle.Parser{}
	default:
		panic(fmt.Sprintf("no parser for series \"%s\"", series))
	}
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
