package cmds

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/got-style/style"
)

type StatCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewStatCommand() *StatCommand {
	return &StatCommand{
		CommandBase: command.NewCommandBase("stat", "calc file stats"),
	}
}

func (cmd *StatCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd StatCommand) Run(args []string) error {
	//ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1

	series := NewSeriesSelect(cmd.Flags)
	cmd.ParseFlags(args)

	url, err := series.Select(PATH)
	if err != nil {
		return err
	}
	style.BoldInfo.Println(url)

	files, err := filepath.Glob(filepath.Join(PATH, url, "src", "v*.mp4"))
	if err != nil {
		return errors.Chain(err, "error reading source directory")
	}

	outFile := filepath.Join(PATH, url, "src", "durations.txt")

	out, err := os.Create(outFile)
	if err != nil {
		return errors.Chain(err, "error creating durations file")
	}
	defer out.Close()

	for _, file := range files {
		err := cmd.calcVideoDuration(file, out)
		if err != nil {
			return errors.Chain(err, "error calculating video duration")
		}
	}

	style.Create.Printf("+ %s\n", outFile)
	return nil
}

func (cmd StatCommand) calcVideoDuration(path string, w io.Writer) error {
	buffer := bytes.Buffer{}

	exe := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", path)
	exe.Stdout = w
	exe.Stderr = &buffer

	if err := exe.Run(); err != nil {
		return errors.Chain(err, buffer.String())
	}

	return nil
}
