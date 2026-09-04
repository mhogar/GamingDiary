package cmds

import (
	"bytes"
	"fmt"
	"gamingdiary/data"
	"gamingdiary/data/heartgold"
	luigi_mansion_videos "gamingdiary/data/luigi_mansion/videos"
	shake_it_videos "gamingdiary/data/shake_it/videos"
	sunshine_chapters "gamingdiary/data/sunshine/chapters"
	sunshine_shorts "gamingdiary/data/sunshine/shorts"
	ttyd_battles "gamingdiary/data/ttyd/battles"
	ttyd_chapters "gamingdiary/data/ttyd/chapters"
	ttyd_shorts "gamingdiary/data/ttyd/shorts"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type Parser interface {
	RawFiles(path string) ([]string, error)
	ParseEntry(path string) (data.Entry, error)
}

type Upgrader interface {
	UpgradeEntry(entry *data.Entry) error
}

//========================================

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
	name := cmd.Flags.String("name", "", "the name of the series")
	public := cmd.Flags.String("public", "public", "the public path")
	upgrade := cmd.Flags.Bool("upgrade", false, "upgrade existing entries")
	cmd.ParseFlags(args)

	if *name == "" {
		return errors.New("\"name\" cannot be empty")
	}
	dataPath := filepath.Join("data", *name)

	series, err := json.UnmarshalFile[data.Series](filepath.Join(dataPath, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading index file")
	}
	style.BoldInfo.Println(*name)

	if *upgrade {
		return cmd.runUpgrade(filepath.Join(dataPath, series.Entries), *name)
	} else {
		return cmd.runBuild(dataPath, *name, series.Entries, *public)
	}
}

func (cmd BuildCommand) runUpgrade(path, name string) error {
	upgrader, err := cmd.selectUpgrader(name)
	if err != nil {
		return err
	}

	entries, err := json.UnmarshalFile[data.Entries](path)
	if err != nil {
		return errors.Chain(err, "error loading entries")
	}

	for i := range entries.Entries {
		if err := upgrader.UpgradeEntry(&entries.Entries[i]); err != nil {
			return errors.Chain(err, "error upgrading entry")
		}
	}

	err = json.MarshalFilePretty(entries, path, "  ")
	if err != nil {
		return errors.Chain(err, "error saving entries")
	}

	style.Success.Printf("Upgraded [%d] entries\n", len(entries.Entries))
	return nil
}

func (cmd BuildCommand) selectUpgrader(series string) (Upgrader, error) {
	switch series {
	case "ttyd/chapters":
		return ttyd_chapters.Upgrader{}, nil
	case "ttyd/battles":
		return ttyd_battles.Upgrader{}, nil
	default:
		return nil, errors.Format("no upgrader for series \"%s\"", series)
	}
}

func (cmd BuildCommand) runBuild(path, name, entires, public string) error {
	parser, err := cmd.selectParser(name)
	if err != nil {
		return err
	}

	files, err := parser.RawFiles(path)
	if err != nil {
		return errors.Chain(err, "error getting raw files")
	}

	entries := data.Entries{
		Entries: make([]data.Entry, len(files)),
	}

	for i, file := range files {
		style.Info.Printf("\r%s ", file)
		entries.VideoCount++

		entry, err := parser.ParseEntry(file)
		if err != nil {
			return errors.Chain(err, "error parsing raw file")
		}

		entry.Duration, err = cmd.calcVideoDuration(filepath.Join(public, name, entry.Video))
		if err != nil {
			style.Error.Printf(" -> [x] %s\n", filepath.Join(name, entry.Video))
		}
		entries.TotalDuration += entry.Duration

		entries.Entries[i] = entry
	}
	fmt.Println()

	out := filepath.Join(path, entires)

	err = json.MarshalFilePretty(entries, out, "  ")
	if err != nil {
		return errors.Chain(err, "error saving entries file")
	}

	style.Create.Printf("+ %s\n", out)
	return nil
}

func (cmd BuildCommand) selectParser(series string) (Parser, error) {
	switch series {
	case "sunshine/chapters":
		return sunshine_chapters.Parser{}, nil
	case "sunshine/shorts":
		return sunshine_shorts.Parser{}, nil
	case "ttyd/chapters":
		return ttyd_chapters.Parser{}, nil
	case "ttyd/battles":
		return ttyd_battles.Parser{}, nil
	case "ttyd/shorts":
		return ttyd_shorts.Parser{}, nil
	case "luigi_mansion/videos":
		return luigi_mansion_videos.Parser{}, nil
	case "shake_it/videos":
		return shake_it_videos.Parser{}, nil
	case "heartgold":
		return heartgold.Parser{}, nil
	default:
		return nil, errors.Format("no parser for series \"%s\"", series)
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
