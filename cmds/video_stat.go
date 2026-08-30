package cmds

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/got-style/style"
)

type VideoStatCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewVideoStatCommand() *VideoStatCommand {
	return &VideoStatCommand{
		CommandBase: command.NewCommandBase("vstat", "calc video stats"),
	}
}

func (cmd *VideoStatCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd VideoStatCommand) Run(args []string) error {
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

	files, err := filepath.Glob(filepath.Join(*src, "v*.mp4"))
	if err != nil {
		return errors.Chain(err, "error reading videos directory")
	}

	if len(files) == 0 {
		return errors.New("no video files found")
	}

	outFile := filepath.Join("data", series, "video_stats.txt")

	out, err := os.Create(outFile)
	if err != nil {
		return errors.Chain(err, "error creating video stats file")
	}
	defer out.Close()

	for _, file := range files {
		fmt.Print(file)

		err := cmd.calcVideoDuration(file, out)
		if err != nil {
			return errors.Chain(err, "error calculating video duration")
		}

		fmt.Print("\r")
	}
	fmt.Println()

	style.Create.Printf("+ %s\n", outFile)
	return nil
}

func (cmd VideoStatCommand) calcVideoDuration(path string, w io.Writer) error {
	buffer := bytes.Buffer{}

	exe := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", path)
	exe.Stdout = w
	exe.Stderr = &buffer

	if err := exe.Run(); err != nil {
		return errors.Chain(err, buffer.String())
	}

	return nil
}
