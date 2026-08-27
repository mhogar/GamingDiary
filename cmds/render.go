package cmds

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type BaseData struct {
	URLs   []string `json:"series"`
	Series []SeriesHeader
}

type SeriesHeader struct {
	URL       string
	Theme     string
	Title     string
	Thumbnail string
	Dates     string
}

type SeriesData struct {
	Theme      string            `json:"theme"`
	Background string            `json:"background"`
	Title      string            `json:"title"`
	Dates      string            `json:"dates"`
	Thumbnail  string            `json:"thumbnail"`
	Groups     map[string]string `json:"groups"`

	ResourcePath string
	Videos       []Video
}

type Video struct {
	Groups      []string
	Title       string
	Description string
	Duration    string
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
	all := cmd.Flags.Bool("all", false, "render all templates")
	cmd.ParseFlags(args)

	if *all || *series.Index == 0 {
		if err := cmd.renderBase(); err != nil {
			return err
		}
	}

	if *all {
		return cmd.renderAllSeries()
	} else if *series.Index != 0 {
		return cmd.renderSingleSeries(series)
	}

	return nil
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
			Theme:     s.Theme,
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

func (cmd RenderCommand) renderAllSeries() error {
	data, err := json.UnmarshalFile[BaseData](filepath.Join(PATH, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	for _, series := range data.URLs {
		if err := cmd.renderSeries(series); err != nil {
			return err
		}
	}

	return nil
}

func (cmd RenderCommand) renderSingleSeries(series SeriesSelect) error {
	s, err := series.Select(PATH)
	if err != nil {
		return err
	}

	style.BoldInfo.Println(s)
	return cmd.renderSeries(s)
}

func (cmd RenderCommand) renderSeries(series string) error {
	t := template.Must(template.ParseFiles("series.gohtml"))

	//-- load data
	data, err := json.UnmarshalFile[SeriesData](filepath.Join(PATH, series, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}
	data.ResourcePath, _ = filepath.Rel(filepath.Join(PATH, series), PATH)

	groups := make(map[string]*regexp.Regexp)
	for key, val := range data.Groups {
		groups[key] = regexp.MustCompile(val)
	}

	durations := []string{}
	bytes, err := os.ReadFile(filepath.Join(PATH, series, "src", "video_stats.txt"))
	if err == nil {
		durations = strings.Split(string(bytes), "\n")
	}

	files, err := filepath.Glob(filepath.Join(PATH, series, "src", "meta*.txt"))
	if err != nil {
		return errors.Chain(err, "error reading source directory")
	}

	data.Videos = make([]Video, len(files))
	for i, file := range files {
		data.Videos[i], _ = cmd.buildVideo(file, groups, cmd.parseDuration(durations, i))
	}

	//-- execute the template
	out := filepath.Join(PATH, series, "index.html")

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

func (RenderCommand) parseDuration(durations []string, index int) int {
	if index >= len(durations) {
		return 0
	}

	f, err := strconv.ParseFloat(durations[index], 32)
	if err != nil {
		return 0
	}

	return int(math.Round(f))
}

func (RenderCommand) buildVideo(path string, groupExps map[string]*regexp.Regexp, duration int) (Video, error) {
	index := regexp.MustCompile(`meta(.+)\.txt$`).FindStringSubmatch(path)[1]

	groups := []string{}
	for group, regex := range groupExps {
		if regex.MatchString(index) {
			groups = append(groups, group)
		}
	}

	meta, err := os.ReadFile(path)
	if err != nil {
		return Video{}, errors.Chain(err, "error reading meta file")
	}
	lines := strings.Split(string(meta), "\n")

	return Video{
		Groups:      groups,
		Title:       fmt.Sprintf("Chapter %s | %s\n", index, strings.SplitN(lines[2], " | ", 2)[0]),
		Description: lines[5],
		Duration:    fmt.Sprintf("%d:%02d:%02d", duration/(60*60), duration/60, duration%60),
		Thumbnail:   fmt.Sprintf("src/t%s.png", index),
		Video:       fmt.Sprintf("src/v%s.mp4", index),
	}, nil
}
