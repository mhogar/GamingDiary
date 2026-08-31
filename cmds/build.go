package cmds

import (
	"bytes"
	"fmt"
	"local/data"
	"local/data/sunshine"
	ttyd_battle "local/data/ttyd/battles"
	ttyd_chapter "local/data/ttyd/chapters"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

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
	name := cmd.Flags.String("name", "", "the name of the series")
	cmd.ParseFlags(args)

	if *public == "" {
		return errors.New("\"public\" cannot be empty")
	}
	if *name == "" {
		return errors.New("\"name\" cannot be empty")
	}

	dataPath := filepath.Join("data", *name)

	series, err := json.UnmarshalFile[data.Series](filepath.Join(dataPath, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading index file")
	}
	style.BoldInfo.Println(dataPath)

	files, err := filepath.Glob(filepath.Join(dataPath, series.Raw))
	if err != nil {
		return errors.Chain(err, "error finding raw files")
	}

	entries := data.Entries{
		Entries: make([]data.Entry, len(files)),
	}

	for i, file := range files {
		style.Info.Printf("\r%s", file)
		entries.VideoCount++

		entry, err := cmd.selectParser(*name).Parse(file)
		if err != nil {
			return errors.Chain(err, "error parsing raw file")
		}

		entry.Duration, err = cmd.calcVideoDuration(filepath.Join(*public, *name, entry.Video))
		if err != nil {
			return errors.Chain(err, "error calculating video duration")
		}
		entries.TotalDuration += entry.Duration

		entries.Entries[i] = entry
	}
	fmt.Println()

	out := filepath.Join(dataPath, series.Entries)

	err = json.MarshalFilePretty(entries, out, "  ")
	if err != nil {
		return errors.Chain(err, "error saving entries file")
	}

	style.Create.Printf("+ %s\n", out)
	return nil
}

func (cmd BuildCommand) selectParser(series string) data.Parser {
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
