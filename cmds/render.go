package cmds

import (
	"app/data"
	"app/data/build"
	"app/tools/melt"
	"fmt"
	"path/filepath"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type RenderCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewRenderCommand() *RenderCommand {
	return &RenderCommand{
		CommandBase: command.NewCommandBase("render", "Render a video"),
	}
}

func (cmd *RenderCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd RenderCommand) Run(args []string) error {
	series := cmd.Flags.String("series", "", "name of the series")
	video := cmd.Flags.String("video", "", "video to render")
	index := cmd.Flags.String("index", "", "entry index to video is for")
	cmd.Flags.Parse(args)

	if *series == "" {
		return errors.New("\"series\" cannot be empty")
	}
	if *video == "" {
		return errors.New("\"video\" cannot be empty")
	}
	if *index == "" {
		return errors.New("\"index\" cannot be empty")
	}

	entryPath := filepath.Join(data.STATIC_PATH, *series, fmt.Sprintf("entry%s.json", *index))
	_, err := json.UnmarshalFile[build.Entry](entryPath)
	if err != nil {
		return errors.Chain(err, "error reading entry")
	}
	style.Bold.Println(entryPath)

	melt.RenderVideo(*video)

	return nil
}
