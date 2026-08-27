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
	GLOB = "meta*.txt"
)

type BaseData struct {
	URLs   []string `json:"series"`
	Series []SeriesHeader
}

type SeriesHeader struct {
	URL       string
	Title     string
	Thumbnail string
	Dates     string
}

type SeriesData struct {
	Title     string            `json:"title"`
	Dates     string            `json:"dates"`
	Thumbnail string            `json:"thumbnail"`
	Groups    map[string]string `json:"groups"`
	Prologue  Video
	Videos    []Video
}

type Video struct {
	Groups      []string
	Title       string
	Description string
	Thumbnail   string
	Video       string
}

//====================================================

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

	if *series.Index == 0 {
		return cmd.renderBase()
	} else {
		return cmd.renderSeries(series)
	}
}

func (cmd RenderCommand) renderBase() error {
	t := template.Must(template.ParseFiles("base.gohtml"))

	data, err := json.UnmarshalFile[BaseData](filepath.Join(PATH, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	data.Series = make([]SeriesHeader, len(data.URLs))
	for i, url := range data.URLs {
		s, err := json.UnmarshalFile[SeriesData](filepath.Join(PATH, url, "index.json"))
		if err != nil {
			return errors.Chain(err, "error reading series data file")
		}

		data.Series[i] = SeriesHeader{
			URL:       url,
			Title:     fmt.Sprintf("(%d) %s", i+1, s.Title),
			Thumbnail: filepath.Join(PATH, url, s.Thumbnail),
			Dates:     s.Dates,
		}
	}

	//-- execute the template
	out := filepath.Join(PATH, "index.html")

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

func (cmd RenderCommand) renderSeries(series SeriesSelect) error {
	t := template.Must(template.ParseFiles("series.gohtml"))

	s, err := series.Select(PATH)
	if err != nil {
		return err
	}
	style.BoldInfo.Println(s)

	//-- load data
	data, err := json.UnmarshalFile[SeriesData](filepath.Join(PATH, s, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	groups := make(map[string]*regexp.Regexp)
	for key, val := range data.Groups {
		groups[key] = regexp.MustCompile(val)
	}

	files, err := filepath.Glob(filepath.Join(PATH, s, "src", GLOB))
	if err != nil {
		return errors.Chain(err, "error reading source directory")
	}

	data.Prologue = cmd.buildVideo(files[0], groups)
	files = files[1:]

	data.Videos = make([]Video, len(files))
	for i, file := range files {
		data.Videos[i] = cmd.buildVideo(file, groups)
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

func (RenderCommand) buildVideo(path string, groupExps map[string]*regexp.Regexp) Video {
	index := regexp.MustCompile(`meta(.+)\.txt$`).FindStringSubmatch(path)[1]

	groups := []string{}
	for group, regex := range groupExps {
		if regex.MatchString(index) {
			groups = append(groups, group)
		}
	}

	meta, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(meta), "\n")

	return Video{
		Groups:      groups,
		Title:       fmt.Sprintf("Chapter %s | %s\n", index, strings.SplitN(lines[2], " | ", 2)[0]),
		Description: lines[5],
		Thumbnail:   fmt.Sprintf("src/t%s.png", index),
		Video:       fmt.Sprintf("src/v%s.mp4", index),
	}
}
