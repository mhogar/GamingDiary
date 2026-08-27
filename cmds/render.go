package cmds

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

const (
	REGEX = `meta(.+)\.txt$`
	GLOB  = "meta*.txt"
)

type Data struct {
	Title    string `json:"title"`
	Dates    string `json:"dates"`
	Prologue Video
	Videos   []Video
}

type Video struct {
	Title       string
	Description string
	Thumbnail   string
	Video       string
}

//=============================================

type RenderCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewRenderCommand() *RenderCommand {
	return &RenderCommand{
		CommandBase: command.NewCommandBase("render", "render the template"),
	}
}

func (cmd *RenderCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd RenderCommand) Run(args []string) error {
	series := NewSeriesSelect(cmd.Flags)
	cmd.ParseFlags(args)

	t := template.Must(template.ParseFiles("template.gohtml"))

	s, err := series.Select(PATH)
	if err != nil {
		return err
	}
	style.BoldInfo.Println(s)

	//-- load data
	data, err := json.UnmarshalFile[Data](filepath.Join(PATH, s, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	files, err := filepath.Glob(filepath.Join(PATH, s, "src", GLOB))
	if err != nil {
		return errors.Chain(err, "error reading source directory")
	}

	data.Prologue = cmd.buildVideo(files[0])
	files = files[1:]

	data.Videos = make([]Video, len(files))
	for i, file := range files {
		data.Videos[i] = cmd.buildVideo(file)
	}

	//-- execute the template
	out := filepath.Join(PATH, s, "index.html")

	file, err := os.Create(out)
	if err != nil {
		return errors.Chain(err, "error creating index file")
	}
	defer file.Close()

	if err := t.Execute(file, data); err != nil {
		return errors.Chain(err, "error executing template")
	}

	style.Create.Printf("+ %s\n", out)
	return nil
}

func (RenderCommand) buildVideo(path string) Video {
	index := regexp.MustCompile(REGEX).FindStringSubmatch(path)[1]

	meta, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(meta), "\n")

	return Video{
		Title:       fmt.Sprintf("Chapter %s | %s\n", index, strings.SplitN(lines[2], " | ", 2)[0]),
		Description: lines[5],
		Thumbnail:   fmt.Sprintf("src/t%s.png", index),
		Video:       fmt.Sprintf("src/v%s.mp4", index),
	}
}
