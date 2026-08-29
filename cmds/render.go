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
	Series  []string `json:"series"`
	Headers []SeriesHeader
}

type SeriesHeader struct {
	URL           string
	Theme         string
	Title         string
	Thumbnail     string
	VideoCount    int
	TotalDuration string
	Dates         string
	Description   string
}

type SeriesData struct {
	Theme       string            `json:"theme"`
	Background  string            `json:"background"`
	Title       string            `json:"title"`
	Dates       string            `json:"dates"`
	Thumbnail   string            `json:"thumbnail"`
	Description string            `json:"description"`
	Groups      map[string]string `json:"groups"`

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
	t := template.Must(template.ParseFiles("templates/base.gohtml"))

	data, err := json.UnmarshalFile[BaseData](filepath.Join(PATH, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	data.Headers = make([]SeriesHeader, len(data.Series))
	for i, series := range data.Series {
		path := filepath.Join(PATH, "public", series)

		s, err := json.UnmarshalFile[SeriesData](filepath.Join(path, "index.json"))
		if err != nil {
			return errors.Chain(err, "error reading series data file")
		}

		durations := cmd.parseDurations(filepath.Join(path, "data", "video_stats.txt"))
		var total float32

		for _, d := range durations {
			total += d
		}

		data.Headers[i] = SeriesHeader{
			URL:           filepath.Join(path, "index.html"),
			Theme:         s.Theme,
			Title:         fmt.Sprintf("(%d) %s", i+1, s.Title),
			Thumbnail:     filepath.Join(series, s.Thumbnail),
			VideoCount:    len(durations),
			TotalDuration: cmd.formatDuration(total),
			Dates:         s.Dates,
			Description:   s.Description,
		}
	}

	//-- execute the template
	out := filepath.Join(PATH, "public", "index.html")

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

	for _, series := range data.Series {
		if err := cmd.renderSeries(filepath.Join(PATH, "public", series)); err != nil {
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
	return cmd.renderSeries(filepath.Join(PATH, "public", s))
}

func (cmd RenderCommand) renderSeries(path string) error {
	t := template.Must(template.ParseFiles("templates/series.gohtml"))

	//-- load data
	data, err := json.UnmarshalFile[SeriesData](filepath.Join(path, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}
	data.ResourcePath, _ = filepath.Rel(path, filepath.Join(PATH, "public"))

	groups := make(map[string]*regexp.Regexp)
	for key, val := range data.Groups {
		groups[key] = regexp.MustCompile(val)
	}

	files, err := filepath.Glob(filepath.Join(path, "data", "meta*.txt"))
	if err != nil {
		return errors.Chain(err, "error reading source directory")
	}

	durations := cmd.parseDurations(filepath.Join(path, "data", "video_stats.txt"))

	data.Videos = make([]Video, len(files))
	for i, file := range files {
		var duration float32
		if i < len(durations) {
			duration = durations[i]
		}

		data.Videos[i], _ = cmd.buildVideo(file, groups, duration)
	}

	//-- execute the template
	out := filepath.Join(path, "index.html")

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

func (RenderCommand) parseDurations(path string) []float32 {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return []float32{}
	}
	strs := strings.Split(string(bytes), "\n")

	durations := make([]float32, len(strs))

	for i, str := range strs {
		f, err := strconv.ParseFloat(str, 32)
		if err == nil {
			durations[i] = float32(f)
		}
	}

	return durations
}

func (cmd RenderCommand) buildVideo(path string, groupExps map[string]*regexp.Regexp, duration float32) (Video, error) {
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
		Duration:    cmd.formatDuration(duration),
		Thumbnail:   fmt.Sprintf("imgs/t%s.png", index),
		Video:       fmt.Sprintf("videos/v%s.mp4", index),
	}, nil
}

func (RenderCommand) formatDuration(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%02d:%02d:%02d", d/(60*60), (d/60)%60, d%60)
}
